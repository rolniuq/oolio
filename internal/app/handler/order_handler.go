package handler

import (
	"net/http"

	"oolio/internal/app/models"
	"oolio/internal/app/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service      services.OrderService
	queueService services.OrderQueueService
}

func NewOrderHandler(service services.OrderService, queueService services.OrderQueueService) *OrderHandler {
	return &OrderHandler{
		service:      service,
		queueService: queueService,
	}
}

func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	ctx := c.Request.Context()

	var orderReq models.OrderReq
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request format",
		})
		return
	}

	// Validate request
	if len(orderReq.Items) == 0 {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Order must contain at least one item",
		})
		return
	}

	// Validate each item
	for _, item := range orderReq.Items {
		// Validate product ID format (UUID)
		if item.ProductID == "" || len(item.ProductID) != 36 {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Code:    http.StatusBadRequest,
				Type:    "error",
				Message: "Invalid product ID format",
			})
			return
		}

		// Validate quantity
		if item.Quantity <= 0 {
			c.JSON(http.StatusUnprocessableEntity, models.ApiResponse{
				Code:    http.StatusUnprocessableEntity,
				Type:    "error",
				Message: "Quantity must be greater than 0",
			})
			return
		}
	}

	// Add order to queue for batch processing
	queueItem, err := h.queueService.AddOrderToQueue(ctx, &orderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Failed to queue order",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":     "Order queued for processing",
		"queueItemId": queueItem.ID,
		"status":      queueItem.Status,
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("orderId")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Order ID is required",
		})
		return
	}

	// First try to get order from queue (for recent orders)
	queueItem, err := h.queueService.GetOrderFromQueue(ctx, orderID)
	if err == nil && queueItem.Order != nil {
		c.JSON(http.StatusOK, queueItem.Order)
		return
	}

	// If not found in queue, try the orders table
	order, err := h.service.GetOrder(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, models.ApiResponse{
				Code:    http.StatusNotFound,
				Type:    "error",
				Message: "Order not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Failed to retrieve order",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all orders from queue
	orders, err := h.queueService.GetCompletedOrders(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Failed to get orders",
		})
		return
	}

	// Get stats for additional information
	stats, err := h.queueService.GetQueueStatus(ctx)
	if err != nil {
		// Continue without stats if it fails
		stats = make(map[string]int)
	}

	// Transform queue items to order display format
	orderList := make([]gin.H, 0)
	for _, item := range orders {
		orderDisplay := gin.H{
			"id":        item.ID,
			"status":    item.Status,
			"createdAt": item.CreatedAt,
			"updatedAt": item.UpdatedAt,
			"customer":  "Guest", // Default customer name
		}

		// Add order data if available
		if item.Order != nil {
			orderDisplay["total"] = item.Order.Total
			if item.Order.Items != nil {
				orderDisplay["items"] = item.Order.Items
			}
		} else {
			// Calculate total from order request if order data not available
			total := 0.0
			if len(item.OrderReq.Items) > 0 {
				items := make([]gin.H, 0)
				for _, reqItem := range item.OrderReq.Items {
					total += reqItem.Price * float64(reqItem.Quantity)
					items = append(items, gin.H{
						"productId": reqItem.ProductID,
						"price":     reqItem.Price,
						"quantity":  reqItem.Quantity,
					})
				}
				orderDisplay["items"] = items
			}
			orderDisplay["total"] = total
		}

		// Add error message if failed
		if item.Status == "failed" && item.Error != "" {
			orderDisplay["error"] = item.Error
		}

		orderList = append(orderList, orderDisplay)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":  orderList,
		"stats":   stats,
		"message": "Orders retrieved successfully",
	})
}

func (h *OrderHandler) GetQueueStatus(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.queueService.GetQueueStatus(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Failed to get queue status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"queueStats": stats,
	})
}
