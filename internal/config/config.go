package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	API      APIConfig
	Coupon   CouponConfig
	Redis    RedisConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Port string
	Host string
}

type APIConfig struct {
	APIKey string
}

type CouponConfig struct {
	BaseURL string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "oolio"),
			Password: getEnv("DB_PASSWORD", "oolio_password"),
			DBName:   getEnv("DB_NAME", "oolio_db"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		API: APIConfig{
			APIKey: getEnv("API_KEY", "apitest"),
		},
		Coupon: CouponConfig{
			BaseURL: getEnv("COUPON_BASE_URL", "https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
	}
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
