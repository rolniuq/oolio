package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// IntegrationTestClient extends the basic test client for integration testing
type IntegrationTestClient struct {
	*TestClient
	baseURL string
	apiKey  string
}

// NewIntegrationTestClient creates a new integration test client
func NewIntegrationTestClient() *IntegrationTestClient {
	return &IntegrationTestClient{
		TestClient: NewTestClient(),
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// TestCompleteOrderWorkflow tests the complete order placement workflow
func TestCompleteOrderWorkflow(t *testing.T) {
	client := NewIntegrationTestClient()
	
	// Step 1: Get available products
	productsResp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer productsResp.Body.Close()
	
	if productsResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", productsResp.StatusCode)
	}
	
	var products []Product
	if err := json.NewDecoder(productsResp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %v", err)
	}
	
	if len(products) == 0 {
		t.Skip("No products available for integration testing")
	}
	
	// Step 2: Get details for a specific product
	selectedProduct := products[0]
	productResp, err := client.makeRequest("GET", fmt.Sprintf("/product/%s", selectedProduct.ID), nil)
	if err != nil {
		t.Fatalf("Failed to get product details: %v", err)
	}
	defer productResp.Body.Close()
	
	if productResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", productResp.StatusCode)
	}
	
	var productDetail Product
	if err := json.NewDecoder(productResp.Body).Decode(&productDetail); err != nil {
		t.Fatalf("Failed to decode product detail: %v", err)
	}
	
	// Verify product details match
	if productDetail.ID != selectedProduct.ID {
		t.Errorf("Expected product ID %s, got %s", selectedProduct.ID, productDetail.ID)
	}
	
	if productDetail.Name != selectedProduct.Name {
		t.Errorf("Expected product name %s, got %s", selectedProduct.Name, productDetail.Name)
	}
	
	// Step 3: Place an order with the selected product
	orderReq := OrderRequest{
		Items: []OrderItem{
			{ProductID: selectedProduct.ID, Quantity: 2},
		},
	}
	
	orderResp, err := client.makeRequest("POST", "/order", orderReq)
	if err != nil {
		t.Fatalf("Failed to place order: %v", err)
	}
	defer orderResp.Body.Close()
	
	if orderResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", orderResp.StatusCode)
	}
	
	var order Order
	if err := json.NewDecoder(orderResp.Body).Decode(&order); err != nil {
		t.Fatalf("Failed to decode order: %v", err)
	}
	
	// Step 4: Verify order details
	if order.ID == "" {
		t.Error("Order ID should not be empty")
	}
	
	if len(order.Items) != 1 {
		t.Errorf("Expected 1 item in order, got %d", len(order.Items))
	}
	
	if order.Items[0].ProductID != selectedProduct.ID {
		t.Errorf("Expected product ID %s, got %s", selectedProduct.ID, order.Items[0].ProductID)
	}
	
	if order.Items[0].Quantity != 2 {
		t.Errorf("Expected quantity 2, got %d", order.Items[0].Quantity)
	}
	
	if len(order.Products) != 1 {
		t.Errorf("Expected 1 product in order, got %d", len(order.Products))
	}
	
	if order.Products[0].ID != selectedProduct.ID {
		t.Errorf("Expected product ID %s, got %s", selectedProduct.ID, order.Products[0].ID)
	}
}

// TestOrderWithPromoCodeWorkflow tests order placement with promo code validation
func TestOrderWithPromoCodeWorkflow(t *testing.T) {
	client := NewIntegrationTestClient()
	
	// Step 1: Get available products
	productsResp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer productsResp.Body.Close()
	
	var products []Product
	if err := json.NewDecoder(productsResp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %v", err)
	}
	
	if len(products) == 0 {
		t.Skip("No products available for integration testing")
	}
	
	// Step 2: Test orders with different promo codes
	testCases := []struct {
		name       string
		couponCode string
		expectSuccess bool
		description string
	}{
		{"Valid promo code", "HAPPYHRS", true, "Should apply discount successfully"},
		{"Invalid promo code", "INVALID123", true, "Order should succeed but coupon not applied"},
		{"No promo code", "", true, "Order should succeed without discount"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderReq := OrderRequest{
				CouponCode: tc.couponCode,
				Items: []OrderItem{
					{ProductID: products[0].ID, Quantity: 1},
				},
			}
			
			orderResp, err := client.makeRequest("POST", "/order", orderReq)
			if err != nil {
				t.Fatalf("Failed to place order: %v", err)
			}
			defer orderResp.Body.Close()
			
			if tc.expectSuccess && orderResp.StatusCode != http.StatusOK {
				t.Errorf("Expected successful order, got status %d: %s", 
					orderResp.StatusCode, tc.description)
			}
			
			if !tc.expectSuccess && orderResp.StatusCode == http.StatusOK {
				t.Errorf("Expected order to fail, but got success: %s", tc.description)
			}
			
			if orderResp.StatusCode == http.StatusOK {
				var order Order
				if err := json.NewDecoder(orderResp.Body).Decode(&order); err != nil {
					t.Fatalf("Failed to decode order: %v", err)
				}
				
				// Verify order structure
				if order.ID == "" {
					t.Error("Order ID should not be empty")
				}
				
				if len(order.Items) == 0 {
					t.Error("Order should contain items")
				}
				
				if len(order.Products) == 0 {
					t.Error("Order should contain product details")
				}
			}
		})
	}
}

// TestMultiItemOrderWorkflow tests orders with multiple items
func TestMultiItemOrderWorkflow(t *testing.T) {
	client := NewIntegrationTestClient()
	
	// Step 1: Get available products
	productsResp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer productsResp.Body.Close()
	
	var products []Product
	if err := json.NewDecoder(productsResp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %v", err)
	}
	
	if len(products) < 2 {
		t.Skip("Need at least 2 products for multi-item order testing")
	}
	
	// Step 2: Create order with multiple items
	orderReq := OrderRequest{
		CouponCode: "FIFTYOFF", // Valid promo code
		Items: []OrderItem{
			{ProductID: products[0].ID, Quantity: 2},
			{ProductID: products[1].ID, Quantity: 1},
		},
	}
	
	orderResp, err := client.makeRequest("POST", "/order", orderReq)
	if err != nil {
		t.Fatalf("Failed to place multi-item order: %v", err)
	}
	defer orderResp.Body.Close()
	
	if orderResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", orderResp.StatusCode)
	}
	
	var order Order
	if err := json.NewDecoder(orderResp.Body).Decode(&order); err != nil {
		t.Fatalf("Failed to decode order: %v", err)
	}
	
	// Step 3: Verify multi-item order details
	if len(order.Items) != 2 {
		t.Errorf("Expected 2 items in order, got %d", len(order.Items))
	}
	
	if len(order.Products) != 2 {
		t.Errorf("Expected 2 products in order, got %d", len(order.Products))
	}
	
	// Verify total quantities match
	totalQuantity := 0
	for _, item := range order.Items {
		totalQuantity += item.Quantity
	}
	
	expectedQuantity := 0
	for _, item := range orderReq.Items {
		expectedQuantity += item.Quantity
	}
	
	if totalQuantity != expectedQuantity {
		t.Errorf("Expected total quantity %d, got %d", expectedQuantity, totalQuantity)
	}
}

// TestErrorHandlingWorkflow tests error handling in complete workflows
func TestErrorHandlingWorkflow(t *testing.T) {
	client := NewIntegrationTestClient()
	
	// Test 1: Get non-existent product
	t.Run("Non-existent product", func(t *testing.T) {
		resp, err := client.makeRequest("GET", "/product/999999", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
	
	// Test 2: Place order with invalid product ID
	t.Run("Invalid product in order", func(t *testing.T) {
		orderReq := OrderRequest{
			Items: []OrderItem{
				{ProductID: "999999", Quantity: 1},
			},
		}
		
		resp, err := client.makeRequest("POST", "/order", orderReq)
		if err != nil {
			t.Fatalf("Failed to place order: %v", err)
		}
		defer resp.Body.Close()
		
		// Should return error for invalid product
		if resp.StatusCode == http.StatusOK {
			t.Error("Expected order to fail with invalid product ID")
		}
	})
	
	// Test 3: Place order with empty items
	t.Run("Empty order items", func(t *testing.T) {
		orderReq := OrderRequest{
			Items: []OrderItem{},
		}
		
		resp, err := client.makeRequest("POST", "/order", orderReq)
		if err != nil {
			t.Fatalf("Failed to place order: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status 422, got %d", resp.StatusCode)
		}
	})
	
	// Test 4: Invalid product ID format
	t.Run("Invalid product ID format", func(t *testing.T) {
		resp, err := client.makeRequest("GET", "/product/abc", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

// TestConcurrentRequests tests concurrent API access
func TestConcurrentRequests(t *testing.T) {
	client := NewIntegrationTestClient()
	
	// Get products first to have valid IDs
	productsResp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer productsResp.Body.Close()
	
	var products []Product
	if err := json.NewDecoder(productsResp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %v", err)
	}
	
	if len(products) == 0 {
		t.Skip("No products available for concurrent testing")
	}
	
	// Test concurrent product requests
	productID := products[0].ID
	concurrentRequests := 10
	results := make(chan error, concurrentRequests)
	
	// Launch concurrent requests
	for i := 0; i < concurrentRequests; i++ {
		go func() {
			resp, err := client.makeRequest("GET", fmt.Sprintf("/product/%s", productID), nil)
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
				return
			}
			
			results <- nil
		}()
	}
	
	// Collect results
	for i := 0; i < concurrentRequests; i++ {
		if err := <-results; err != nil {
			t.Errorf("Concurrent request failed: %v", err)
		}
	}
}

// TestAPISessionWorkflow tests maintaining session state across multiple requests
func TestAPISessionWorkflow(t *testing.T) {
	client := NewIntegrationTestClient()
	
	// Step 1: Authenticate and get products
	productsResp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer productsResp.Body.Close()
	
	if productsResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", productsResp.StatusCode)
	}
	
	var products []Product
	if err := json.NewDecoder(productsResp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %v", err)
	}
	
	// Step 2: Place multiple orders in the same session
	if len(products) > 0 {
		productID := products[0].ID
		
		for i := 0; i < 3; i++ {
			orderReq := OrderRequest{
				Items: []OrderItem{
					{ProductID: productID, Quantity: 1},
				},
			}
			
			orderResp, err := client.makeRequest("POST", "/order", orderReq)
			if err != nil {
				t.Fatalf("Failed to place order %d: %v", i+1, err)
			}
			defer orderResp.Body.Close()
			
			if orderResp.StatusCode != http.StatusOK {
				t.Errorf("Order %d failed with status %d", i+1, orderResp.StatusCode)
			}
		}
	}
}

// TestDataConsistencyWorkflow tests data consistency across API responses
func TestDataConsistencyWorkflow(t *testing.T) {
	client := NewIntegrationTestClient()
	
	// Step 1: Get product list
	productsResp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer productsResp.Body.Close()
	
	var products []Product
	if err := json.NewDecoder(productsResp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %v", err)
	}
	
	if len(products) == 0 {
		t.Skip("No products available for consistency testing")
	}
	
	// Step 2: Get individual product details and verify consistency
	for _, product := range products {
		productResp, err := client.makeRequest("GET", fmt.Sprintf("/product/%s", product.ID), nil)
		if err != nil {
			t.Fatalf("Failed to get product %s: %v", product.ID, err)
		}
		defer productResp.Body.Close()
		
		if productResp.StatusCode != http.StatusOK {
			t.Errorf("Failed to get product %s, status: %d", product.ID, productResp.StatusCode)
			continue
		}
		
		var productDetail Product
		if err := json.NewDecoder(productResp.Body).Decode(&productDetail); err != nil {
			t.Fatalf("Failed to decode product %s: %v", product.ID, err)
		}
		
		// Verify data consistency
		if productDetail.ID != product.ID {
			t.Errorf("Product ID mismatch: list=%s, detail=%s", product.ID, productDetail.ID)
		}
		
		if productDetail.Name != product.Name {
			t.Errorf("Product name mismatch for %s: list=%s, detail=%s", 
				product.ID, product.Name, productDetail.Name)
		}
		
		if productDetail.Price != product.Price {
			t.Errorf("Product price mismatch for %s: list=%f, detail=%f", 
				product.ID, product.Price, productDetail.Price)
		}
	}
}