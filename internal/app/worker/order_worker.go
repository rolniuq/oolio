package worker

import (
	"context"
	"log"
	"time"

	"oolio/internal/app/services"
)

type OrderWorker struct {
	queueService services.OrderQueueService
	interval     time.Duration
	batchSize    int
}

func NewOrderWorker(queueService services.OrderQueueService, interval time.Duration, batchSize int) *OrderWorker {
	return &OrderWorker{
		queueService: queueService,
		interval:     interval,
		batchSize:    batchSize,
	}
}

func (w *OrderWorker) Start(ctx context.Context) {
	log.Printf("Starting order worker with interval %v and batch size %d", w.interval, w.batchSize)
	w.queueService.StartWorker(ctx, w.interval, w.batchSize)
}

func (w *OrderWorker) ProcessBatch(ctx context.Context) error {
	result, err := w.queueService.ProcessBatch(ctx, w.batchSize)
	if err != nil {
		return err
	}

	if result.Processed > 0 || result.Failed > 0 {
		log.Printf("Batch processed: %d succeeded, %d failed", result.Processed, result.Failed)
		if result.Failed > 0 {
			for _, errorMsg := range result.Errors {
				log.Printf("Error: %s", errorMsg)
			}
		}
	}

	return nil
}
