package handler

import (
	"net/http"
	"strings"

	"oolio/internal/app/models"
	"oolio/internal/app/services"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	ctx := c.Request.Context()

	products, err := h.service.GetAllProducts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()
	productID := c.Param("productId")

	if productID == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Product ID is required",
		})
		return
	}

	product, err := h.service.GetProductByID(ctx, productID)
	if err != nil {
		if err != nil && (strings.Contains(err.Error(), "invalid product ID") || strings.Contains(err.Error(), "invalid UUID")) {
			c.JSON(http.StatusBadRequest, models.ApiResponse{
				Code:    http.StatusBadRequest,
				Type:    "error",
				Message: "Invalid product ID format",
			})
			return
		}

		// Check for not found errors
		if err != nil && (strings.Contains(err.Error(), "product not found") || strings.Contains(err.Error(), "failed to get product")) {
			c.JSON(http.StatusNotFound, models.ApiResponse{
				Code:    http.StatusNotFound,
				Type:    "error",
				Message: "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Failed to retrieve product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}
