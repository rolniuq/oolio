package middleware

import (
	"net/http"
	"strings"

	"oolio/internal/app/models"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		for _, err := range c.Errors {
			handleError(c, err)
		}
	}
}

func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle validation errors specifically
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if isValidationError(err.Err) {
					c.JSON(http.StatusBadRequest, models.ApiResponse{
						Code:    http.StatusBadRequest,
						Type:    "validation_error",
						Message: getValidationErrorMessage(err.Err),
					})
					c.Abort()
					return
				}
			}
		}
	}
}

func handleError(c *gin.Context, ginErr *gin.Error) {
	err := ginErr.Err

	// Handle different types of errors
	switch {
	case isValidationError(err):
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "validation_error",
			Message: getValidationErrorMessage(err),
		})
	case isNotFoundError(err):
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Code:    http.StatusNotFound,
			Type:    "not_found",
			Message: "Resource not found",
		})
	case isUnauthorizedError(err):
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Code:    http.StatusUnauthorized,
			Type:    "unauthorized",
			Message: "Unauthorized access",
		})
	case isForbiddenError(err):
		c.JSON(http.StatusForbidden, models.ApiResponse{
			Code:    http.StatusForbidden,
			Type:    "forbidden",
			Message: "Access forbidden",
		})
	case isConflictError(err):
		c.JSON(http.StatusConflict, models.ApiResponse{
			Code:    http.StatusConflict,
			Type:    "conflict",
			Message: "Resource conflict",
		})
	case isUnprocessableEntityError(err):
		c.JSON(http.StatusUnprocessableEntity, models.ApiResponse{
			Code:    http.StatusUnprocessableEntity,
			Type:    "unprocessable_entity",
			Message: "Unprocessable entity",
		})
	default:
		// Log internal errors but don't expose details to client
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "internal_error",
			Message: "Internal server error",
		})
	}
}

func isValidationError(err error) bool {
	return strings.Contains(err.Error(), "validation") ||
		strings.Contains(err.Error(), "required") ||
		strings.Contains(err.Error(), "invalid") ||
		strings.Contains(err.Error(), "must")
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "not exist")
}

func isUnauthorizedError(err error) bool {
	return strings.Contains(err.Error(), "unauthorized") ||
		strings.Contains(err.Error(), "authentication")
}

func isForbiddenError(err error) bool {
	return strings.Contains(err.Error(), "forbidden") ||
		strings.Contains(err.Error(), "permission")
}

func isConflictError(err error) bool {
	return strings.Contains(err.Error(), "conflict") ||
		strings.Contains(err.Error(), "duplicate")
}

func isUnprocessableEntityError(err error) bool {
	return strings.Contains(err.Error(), "invalid coupon") ||
		strings.Contains(err.Error(), "unprocessable")
}

func getValidationErrorMessage(err error) string {
	errStr := err.Error()

	switch {
	case strings.Contains(errStr, "product ID"):
		return "Invalid product ID format"
	case strings.Contains(errStr, "order must contain"):
		return "Order must contain at least one item"
	case strings.Contains(errStr, "quantity"):
		return "Invalid quantity - must be greater than 0"
	case strings.Contains(errStr, "price"):
		return "Invalid price - must be greater than 0"
	case strings.Contains(errStr, "required"):
		return "Required field is missing"
	case strings.Contains(errStr, "name"):
		return "Invalid name format"
	case strings.Contains(errStr, "category"):
		return "Invalid category format"
	default:
		return "Invalid request format"
	}
}

// Recovery middleware for handling panics
func PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, models.ApiResponse{
					Code:    http.StatusInternalServerError,
					Type:    "panic",
					Message: "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
