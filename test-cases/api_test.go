package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

const (
	baseURL = "http://localhost:8082/api/v1"
	apiKey  = "apitest"
)

// TestClient represents a test HTTP client
type TestClient struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// NewTestClient creates a new test client
func NewTestClient() *TestClient {
	return &TestClient{
		client:  &http.Client{},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// makeRequest makes an HTTP request with authentication
func (tc *TestClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := tc.baseURL + endpoint
	
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", tc.apiKey)
	
	return tc.client.Do(req)
}

// Product represents a product entity
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// OrderRequest represents an order request
type OrderRequest struct {
	CouponCode string       `json:"couponCode,omitempty"`
	Items      []OrderItem  `json:"items"`
}

// Order represents an order response
type Order struct {
	ID       string    `json:"id"`
	Items    []OrderItem `json:"items"`
	Products []Product `json:"products"`
}

// TestListProducts tests the GET /product endpoint
func TestListProducts(t *testing.T) {
	client := NewTestClient()
	
	resp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if len(products) == 0 {
		t.Error("Expected at least one product, got none")
	}
	
	// Validate product structure
	for _, product := range products {
		if product.ID == "" {
			t.Error("Product ID is required")
		}
		if product.Name == "" {
			t.Error("Product name is required")
		}
		if product.Price <= 0 {
			t.Error("Product price must be positive")
		}
	}
}

// TestGetProduct tests the GET /product/{productId} endpoint
func TestGetProduct(t *testing.T) {
	client := NewTestClient()
	
	// First get all products to find a valid ID
	listResp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer listResp.Body.Close()
	
	var products []Product
	if err := json.NewDecoder(listResp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode products: %v", err)
	}
	
	if len(products) == 0 {
		t.Skip("No products available for testing")
	}
	
	validProductID := products[0].ID
	
	tests := []struct {
		name       string
		productID  string
		expectCode int
		expectErr  bool
	}{
		{"Valid product ID", validProductID, http.StatusOK, false},
		{"Invalid product ID", "999999", http.StatusBadRequest, false},
		{"Non-numeric product ID", "abc", http.StatusBadRequest, false},
		{"Empty product ID", "", http.StatusOK, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := fmt.Sprintf("/product/%s", tt.productID)
			resp, err := client.makeRequest("GET", endpoint, nil)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, resp.StatusCode)
			}
			
			if !tt.expectErr && resp.StatusCode == http.StatusOK {
				// For empty product ID, we expect an array (redirect to list)
				if tt.productID == "" {
					var products []Product
					if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
						t.Fatalf("Failed to decode response: %v", err)
					}
				} else {
					var product Product
					if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
						t.Fatalf("Failed to decode response: %v", err)
					}
					
					if product.ID != tt.productID {
						t.Errorf("Expected product ID %s, got %s", tt.productID, product.ID)
					}
				}
			}
		})
	}
}

// TestPlaceOrder tests the POST /order endpoint
func TestPlaceOrder(t *testing.T) {
	client := NewTestClient()
	
	// First get available products
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
		t.Skip("No products available for testing")
	}
	
	validProductID := products[0].ID
	
	tests := []struct {
		name       string
		orderReq   OrderRequest
		expectCode int
		expectErr  bool
	}{
		{
			"Valid order without coupon",
			OrderRequest{
				Items: []OrderItem{{ProductID: validProductID, Quantity: 1}},
			},
			http.StatusAccepted,
			false,
		},
		{
			"Valid order with valid coupon",
			OrderRequest{
				CouponCode: "HAPPYHRS",
				Items:      []OrderItem{{ProductID: validProductID, Quantity: 1}},
			},
			http.StatusAccepted,
			false,
		},
		{
			"Valid order with invalid coupon",
			OrderRequest{
				CouponCode: "INVALID",
				Items:      []OrderItem{{ProductID: validProductID, Quantity: 1}},
			},
			http.StatusAccepted, // Order should still be placed, coupon just won't be applied
			false,
		},
		{
			"Empty items array",
			OrderRequest{
				Items: []OrderItem{},
			},
			http.StatusBadRequest,
			true,
		},
		{
			"Invalid product ID",
			OrderRequest{
				Items: []OrderItem{{ProductID: "invalid-uuid-format", Quantity: 1}},
			},
			http.StatusBadRequest,
			true,
		},
		{
			"Zero quantity",
			OrderRequest{
				Items: []OrderItem{{ProductID: validProductID, Quantity: 0}},
			},
			http.StatusUnprocessableEntity,
			true,
		},
		{
			"Negative quantity",
			OrderRequest{
				Items: []OrderItem{{ProductID: validProductID, Quantity: -1}},
			},
			http.StatusUnprocessableEntity,
			true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.makeRequest("POST", "/order", tt.orderReq)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, resp.StatusCode)
			}
			
			if !tt.expectErr && resp.StatusCode == http.StatusOK {
				var order Order
				if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				
				if order.ID == "" {
					t.Error("Order ID is required")
				}
				
				if len(order.Items) == 0 {
					t.Error("Order must have at least one item")
				}
				
				if len(order.Products) == 0 {
					t.Error("Order must include product details")
				}
			}
		})
	}
}

// TestAuthentication tests API key authentication
func TestAuthentication(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		expectCode int
	}{
		{"Valid API key", apiKey, http.StatusOK},
		{"Invalid API key", "invalid_key", http.StatusUnauthorized},
		{"Missing API key", "", http.StatusUnauthorized},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new client with specific API key
			testClient := &TestClient{
				client:  &http.Client{},
				baseURL: baseURL,
				apiKey:  tt.apiKey,
			}
			
			resp, err := testClient.makeRequest("GET", "/product", nil)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, resp.StatusCode)
			}
		})
	}
}

// TestOpenAPICompliance tests compliance with OpenAPI specification
func TestOpenAPICompliance(t *testing.T) {
	client := NewTestClient()
	
	// Test content-type headers
	resp, err := client.makeRequest("GET", "/product", nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
	
	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(actualContentType, expectedContentType) {
		t.Errorf("Expected Content-Type to start with %s, got %s", expectedContentType, actualContentType)
	}
	
	// Test response schema compliance
	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Validate each product matches OpenAPI schema
	for _, product := range products {
		if product.ID == "" {
			t.Error("Product ID is required by OpenAPI schema")
		}
		if product.Name == "" {
			t.Error("Product name is required by OpenAPI schema")
		}
		if product.Price == 0 {
			t.Error("Product price is required by OpenAPI schema")
		}
		// Category is optional in OpenAPI schema, so we don't validate it
	}
}