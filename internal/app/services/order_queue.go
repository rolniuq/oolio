package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"oolio/internal/app/models"
	"oolio/internal/app/repository"

	"github.com/google/uuid"
)

type OrderQueueService interface {
	AddOrderToQueue(ctx context.Context, orderReq *models.OrderReq) (*models.OrderQueueItem, error)
	ProcessBatch(ctx context.Context, batchSize int) (*models.BatchProcessResult, error)
	GetQueueStatus(ctx context.Context) (map[string]int, error)
	GetCompletedOrders(ctx context.Context) ([]*models.OrderQueueItem, error)
	StartWorker(ctx context.Context, interval time.Duration, batchSize int)
	GetOrderFromQueue(ctx context.Context, itemID string) (*models.OrderQueueItem, error)
}

type orderQueueService struct {
	queueRepo repository.OrderQueueRepository
	orderRepo repository.OrderRepository
	orderSvc  OrderService
}

func NewOrderQueueService(queueRepo repository.OrderQueueRepository, orderRepo repository.OrderRepository, orderSvc OrderService) OrderQueueService {
	return &orderQueueService{
		queueRepo: queueRepo,
		orderRepo: orderRepo,
		orderSvc:  orderSvc,
	}
}

func (s *orderQueueService) AddOrderToQueue(ctx context.Context, orderReq *models.OrderReq) (*models.OrderQueueItem, error) {
	item := &models.OrderQueueItem{
		ID:         generateUUID(),
		OrderReq:   *orderReq,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		RetryCount: 0,
	}

	if err := s.queueRepo.AddToQueue(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to add order to queue: %w", err)
	}

	return item, nil
}

func (s *orderQueueService) ProcessBatch(ctx context.Context, batchSize int) (*models.BatchProcessResult, error) {
	items, err := s.queueRepo.GetPendingItems(ctx, batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending items: %w", err)
	}

	if len(items) == 0 {
		return &models.BatchProcessResult{
			Processed: 0,
			Failed:    0,
			Errors:    []string{},
			Items:     []models.OrderQueueItem{},
		}, nil
	}

	result := &models.BatchProcessResult{
		Processed: 0,
		Failed:    0,
		Errors:    []string{},
		Items:     make([]models.OrderQueueItem, 0, len(items)),
	}

	for _, item := range items {
		if err := s.processQueueItem(ctx, item); err != nil {
			log.Printf("Failed to process queue item %s: %v", item.ID, err)
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Item %s: %v", item.ID, err))
		} else {
			result.Processed++
		}
		result.Items = append(result.Items, *item)
	}

	return result, nil
}

func (s *orderQueueService) processQueueItem(ctx context.Context, item *models.OrderQueueItem) error {
	item.Status = "processing"
	item.UpdatedAt = time.Now()

	if err := s.queueRepo.UpdateItem(ctx, item); err != nil {
		return fmt.Errorf("failed to mark item as processing: %w", err)
	}

	order, err := s.orderSvc.CreateOrder(ctx, &item.OrderReq)
	if err != nil {
		item.Status = "failed"
		item.Error = err.Error()
		item.UpdatedAt = time.Now()
		item.RetryCount++

		if item.RetryCount >= 3 {
			log.Printf("Item %s exceeded max retry count, marking as permanently failed", item.ID)
		}

		if updateErr := s.queueRepo.UpdateItem(ctx, item); updateErr != nil {
			return fmt.Errorf("failed to mark item as failed: %w (original error: %v)", updateErr, err)
		}
		return fmt.Errorf("failed to create order: %w", err)
	}

	item.Status = "completed"
	item.Order = order
	item.Error = ""
	item.UpdatedAt = time.Now()

	if err := s.queueRepo.MarkAsCompleted(ctx, item.ID, order); err != nil {
		return fmt.Errorf("failed to mark item as completed: %w", err)
	}

	return nil
}

func (s *orderQueueService) GetQueueStatus(ctx context.Context) (map[string]int, error) {
	return s.queueRepo.GetQueueStats(ctx)
}

func (s *orderQueueService) StartWorker(ctx context.Context, interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Starting order queue worker with interval %v and batch size %d", interval, batchSize)

	for {
		select {
		case <-ctx.Done():
			log.Println("Order queue worker stopped")
			return
		case <-ticker.C:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Worker panic recovered: %v", r)
					}
				}()

				result, err := s.ProcessBatch(ctx, batchSize)
				if err != nil {
					log.Printf("Failed to process batch: %v", err)
					return
				}

				if result.Processed > 0 || result.Failed > 0 {
					log.Printf("Batch processed: %d succeeded, %d failed", result.Processed, result.Failed)
					if result.Failed > 0 {
						for _, errorMsg := range result.Errors {
							log.Printf("Error: %s", errorMsg)
						}
					}
				}
			}()
		}
	}
}

func (s *orderQueueService) GetCompletedOrders(ctx context.Context) ([]*models.OrderQueueItem, error) {
	return s.queueRepo.GetAllOrders(ctx)
}

func (s *orderQueueService) GetOrderFromQueue(ctx context.Context, itemID string) (*models.OrderQueueItem, error) {
	return s.queueRepo.GetOrderFromQueue(ctx, itemID)
}

func generateUUID() string {
	return uuid.New().String()
}
