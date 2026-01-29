package models

import "time"

type OrderItem struct {
	ProductID string  `json:"productId" description:"ID of the product"`
	Quantity  int     `json:"quantity" description:"Item count"`
	Price     float64 `json:"price" description:"Price at time of order"`
}

type Order struct {
	ID        string      `json:"id" example:"0000-0000-0000-0000"`
	Total     float64     `json:"total" example:"90.0"`
	Discounts float64     `json:"discounts" example:"10.0"`
	Items     []OrderItem `json:"items"`
	Products  []Product   `json:"products"`
}

type OrderQueueItem struct {
	ID         string    `json:"id"`
	OrderReq   OrderReq  `json:"orderReq"`
	Status     string    `json:"status"` // pending, processing, completed, failed
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Error      string    `json:"error,omitempty"`
	Order      *Order    `json:"order,omitempty"`
	RetryCount int       `json:"retryCount"`
}

type BatchProcessResult struct {
	Processed int              `json:"processed"`
	Failed    int              `json:"failed"`
	Errors    []string         `json:"errors"`
	Items     []OrderQueueItem `json:"items"`
}
