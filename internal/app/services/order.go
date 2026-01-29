package services

import (
	"context"
	"fmt"

	"oolio/internal/app/models"
	"oolio/internal/app/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, orderReq *models.OrderReq) (*models.Order, error)
	GetOrder(ctx context.Context, id string) (*models.Order, error)
}

type orderService struct {
	orderRepo     repository.OrderRepository
	productRepo   repository.ProductRepository
	couponService CouponService
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, couponService CouponService) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		couponService: couponService,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, orderReq *models.OrderReq) (*models.Order, error) {
	if err := s.validateOrderReq(orderReq); err != nil {
		return nil, fmt.Errorf("order validation failed: %w", err)
	}

	// Get products for all items in the order
	productIDs := make([]string, len(orderReq.Items))
	for i, item := range orderReq.Items {
		productIDs[i] = item.ProductID
	}

	products, err := s.getProductsForOrder(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products for order: %w", err)
	}

	// Calculate order total
	total, err := s.calculateOrderTotal(orderReq.Items, products)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate order total: %w", err)
	}

	// Apply discount if coupon code provided
	discounts := 0.0
	if orderReq.CouponCode != "" {
		discounts, err = s.applyDiscount(total, orderReq.CouponCode)
		if err != nil {
			return nil, fmt.Errorf("failed to apply discount: %w", err)
		}
	}

	// Create order
	order := &models.Order{
		Total:     total,
		Discounts: discounts,
		Items:     orderReq.Items,
		Products:  products,
	}

	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}

	order, err := s.orderRepo.FindOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by ID %s: %w", id, err)
	}

	return order, nil
}

func (s *orderService) validateOrderReq(orderReq *models.OrderReq) error {
	if orderReq == nil {
		return fmt.Errorf("order request cannot be nil")
	}

	if len(orderReq.Items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}

	for i, item := range orderReq.Items {
		if item.ProductID == "" {
			return fmt.Errorf("item %d: product ID is required", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity must be greater than 0", i)
		}
	}

	return nil
}

func (s *orderService) getProductsForOrder(ctx context.Context, productIDs []string) ([]models.Product, error) {
	products := make([]models.Product, 0, len(productIDs))

	for _, productID := range productIDs {
		product, err := s.productRepo.FindOne(ctx, productID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %s: %w", productID, err)
		}
		products = append(products, *product)
	}

	return products, nil
}

func (s *orderService) calculateOrderTotal(items []models.OrderItem, products []models.Product) (float64, error) {
	productPrices := make(map[string]float64)
	for _, product := range products {
		productPrices[product.ID] = product.Price
	}

	total := 0.0
	for _, item := range items {
		price, exists := productPrices[item.ProductID]
		if !exists {
			return 0, fmt.Errorf("product %s not found in order items", item.ProductID)
		}

		itemTotal := price * float64(item.Quantity)
		total += itemTotal
	}

	return total, nil
}

func (s *orderService) applyDiscount(total float64, couponCode string) (float64, error) {
	if !s.couponService.ValidateCoupon(couponCode) {
		return 0, fmt.Errorf("invalid coupon code: %s", couponCode)
	}

	discountPercentage := s.couponService.GetDiscountPercentage(couponCode)
	if discountPercentage <= 0 || discountPercentage > 100 {
		return 0, fmt.Errorf("invalid discount percentage: %f", discountPercentage)
	}

	discount := (total * discountPercentage) / 100
	return discount, nil
}
