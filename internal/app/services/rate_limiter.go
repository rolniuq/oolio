package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiterService interface {
	AllowRequest(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
	IsAllowed(ctx context.Context, key string) (bool, error)
	GetRemainingTokens(ctx context.Context, key string, limit int) (int, error)
	ResetKey(ctx context.Context, key string) error
}

type rateLimiterService struct {
	redisClient *redis.Client
	luaScript   *redis.Script
}

// Lua script for token bucket algorithm
const tokenBucketScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local tokens = tonumber(ARGV[2])
local interval = tonumber(ARGV[3])
local request = tonumber(ARGV[4])
local current_time = tonumber(ARGV[5])

local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
local current_tokens = tonumber(bucket[1]) or capacity
local last_refill = tonumber(bucket[2]) or current_time

local elapsed = current_time - last_refill
local tokens_to_add = math.floor(elapsed / interval) * tokens
current_tokens = math.min(capacity, current_tokens + tokens_to_add)

if current_tokens >= request then
    current_tokens = current_tokens - request
    redis.call('HMSET', key, 'tokens', current_tokens, 'last_refill', current_time)
    redis.call('EXPIRE', key, math.ceil(capacity / tokens * interval))
    return {1, current_tokens}
else
    redis.call('HMSET', key, 'tokens', current_tokens, 'last_refill', current_time)
    redis.call('EXPIRE', key, math.ceil(capacity / tokens * interval))
    return {0, current_tokens}
end
`

func NewRateLimiterService(redisAddr, redisPassword string, redisDB int) RateLimiterService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	script := redis.NewScript(tokenBucketScript)

	return &rateLimiterService{
		redisClient: rdb,
		luaScript:   script,
	}
}

func (s *rateLimiterService) AllowRequest(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	// Calculate token refill rate
	tokensPerWindow := float64(limit)
	intervalMs := window.Milliseconds()

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)

	result, err := s.luaScript.Run(ctx, s.redisClient, []string{key},
		limit,           // capacity
		tokensPerWindow, // tokens per refill
		intervalMs,      // interval in ms
		1,               // tokens requested (1 request)
		currentTime,     // current timestamp in ms
	).Result()

	if err != nil {
		return false, fmt.Errorf("failed to execute rate limiter script: %w", err)
	}

	resultSlice, ok := result.([]interface{})
	if !ok || len(resultSlice) < 2 {
		return false, fmt.Errorf("unexpected result format from rate limiter script")
	}

	allowed, ok := resultSlice[0].(int64)
	if !ok {
		return false, fmt.Errorf("failed to parse allowed value from rate limiter result")
	}

	return allowed == 1, nil
}

func (s *rateLimiterService) IsAllowed(ctx context.Context, key string) (bool, error) {
	return s.AllowRequest(ctx, key, 100, time.Minute) // Default: 100 requests per minute
}

func (s *rateLimiterService) GetRemainingTokens(ctx context.Context, key string, limit int) (int, error) {
	result, err := s.redisClient.HGet(ctx, key, "tokens").Result()
	if err == redis.Nil {
		return limit, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get remaining tokens: %w", err)
	}

	tokens, err := strconv.Atoi(result)
	if err != nil {
		return 0, fmt.Errorf("failed to parse tokens: %w", err)
	}

	return tokens, nil
}

func (s *rateLimiterService) ResetKey(ctx context.Context, key string) error {
	err := s.redisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to reset rate limiter key: %w", err)
	}
	return nil
}
