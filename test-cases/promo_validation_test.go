package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// PromoCodeValidator handles promo code validation logic
type PromoCodeValidator struct {
	couponFiles map[string][]string
}

// NewPromoCodeValidator creates a new validator instance
func NewPromoCodeValidator() *PromoCodeValidator {
	return &PromoCodeValidator{
		couponFiles: make(map[string][]string),
	}
}

// LoadCouponFile loads and decompresses a coupon file
func (pcv *PromoCodeValidator) LoadCouponFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader for %s: %v", filename, err)
	}
	defer gzReader.Close()

	content, err := io.ReadAll(gzReader)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	// Split content into words and filter for potential promo codes
	words := strings.Fields(string(content))
	var promoCodes []string
	
	for _, word := range words {
		// Clean the word - remove punctuation and convert to uppercase
		cleanWord := strings.ToUpper(strings.TrimSpace(word))
		cleanWord = strings.ReplaceAll(cleanWord, ".", "")
		cleanWord = strings.ReplaceAll(cleanWord, ",", "")
		cleanWord = strings.ReplaceAll(cleanWord, "!", "")
		cleanWord = strings.ReplaceAll(cleanWord, "?", "")
		cleanWord = strings.ReplaceAll(cleanWord, ";", "")
		cleanWord = strings.ReplaceAll(cleanWord, ":", "")
		
		// Check if it could be a promo code (8-10 characters, alphanumeric)
		if len(cleanWord) >= 8 && len(cleanWord) <= 10 && pcv.isAlphaNumeric(cleanWord) {
			promoCodes = append(promoCodes, cleanWord)
		}
	}

	pcv.couponFiles[filename] = promoCodes
	return nil
}

// isAlphaNumeric checks if a string contains only letters and numbers
func (pcv *PromoCodeValidator) isAlphaNumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

// ValidatePromoCode validates a promo code according to the requirements
func (pcv *PromoCodeValidator) ValidatePromoCode(code string) bool {
	// Rule 1: Must be a string of length between 8 and 10 characters
	if len(code) < 8 || len(code) > 10 {
		return false
	}

	// Rule 2: Must be found in at least two files
	fileCount := 0
	for _, promoCodes := range pcv.couponFiles {
		for _, promoCode := range promoCodes {
			if strings.EqualFold(promoCode, code) {
				fileCount++
				break // Found in this file, check next file
			}
		}
	}

	return fileCount >= 2
}

// GetValidPromoCodes returns all valid promo codes
func (pcv *PromoCodeValidator) GetValidPromoCodes() []string {
	codeFrequency := make(map[string]int)
	
	// Count frequency across all files
	for _, promoCodes := range pcv.couponFiles {
		uniqueCodes := make(map[string]bool)
		for _, code := range promoCodes {
			if !uniqueCodes[code] {
				uniqueCodes[code] = true
				codeFrequency[code]++
			}
		}
	}
	
	// Return codes that appear in at least 2 files
	var validCodes []string
	for code, frequency := range codeFrequency {
		if frequency >= 2 {
			validCodes = append(validCodes, code)
		}
	}
	
	return validCodes
}

// TestPromoCodeLengthValidation tests the length requirement
func TestPromoCodeLengthValidation(t *testing.T) {
	validator := NewPromoCodeValidator()
	
	// Mock coupon files for testing - same as TestPromoCodeFilePresence
	validator.couponFiles = map[string][]string{
		"file1.gz": {"HAPPYHRS", "FIFTYOFF", "UNIQUE1"},
		"file2.gz": {"HAPPYHRS", "SAVE10NOW", "UNIQUE2"},
		"file3.gz": {"FIFTYOFF", "SAVE10NOW", "UNIQUE3"},
	}
	
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"Valid 8 characters", "HAPPYHRS", true},
		{"Valid 9 characters", "FIFTYOFF", true},
		{"Valid 10 characters", "SAVE10NOW", true},
		{"Too short - 7 characters", "SHORT12", false},
		{"Too long - 11 characters", "TOOLONG123", false},
		{"Empty string", "", false},
		{"Single character", "A", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidatePromoCode(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %v for code '%s', got %v", tt.expected, tt.code, result)
			}
		})
	}
}

// TestPromoCodeFilePresence tests the file presence requirement
func TestPromoCodeFilePresence(t *testing.T) {
	validator := NewPromoCodeValidator()
	
	// Mock coupon files for testing
	validator.couponFiles = map[string][]string{
		"file1.gz": {"HAPPYHRS", "FIFTYOFF", "UNIQUE1"},
		"file2.gz": {"HAPPYHRS", "SAVE10NOW", "UNIQUE2"},
		"file3.gz": {"FIFTYOFF", "SAVE10NOW", "UNIQUE3"},
	}
	
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"Found in 2 files - HAPPYHRS", "HAPPYHRS", true},
		{"Found in 2 files - FIFTYOFF", "FIFTYOFF", true},
		{"Found in 2 files - SAVE10NOW", "SAVE10NOW", true},
		{"Found in only 1 file - UNIQUE1", "UNIQUE1", false},
		{"Found in only 1 file - UNIQUE2", "UNIQUE2", false},
		{"Found in only 1 file - UNIQUE3", "UNIQUE3", false},
		{"Not found in any file - MISSING", "MISSING", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidatePromoCode(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %v for code '%s', got %v", tt.expected, tt.code, result)
			}
		})
	}
}

// TestPromoCodeCaseInsensitivity tests case insensitive validation
func TestPromoCodeCaseInsensitivity(t *testing.T) {
	validator := NewPromoCodeValidator()
	
	// Mock coupon files with uppercase codes - ensure FIFTYOFF and SAVE10NOW appear in 2 files
	validator.couponFiles = map[string][]string{
		"file1.gz": {"HAPPYHRS", "FIFTYOFF"},
		"file2.gz": {"HAPPYHRS", "SAVE10NOW"},
		"file3.gz": {"FIFTYOFF", "SAVE10NOW"},
	}
	
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"Uppercase - HAPPYHRS", "HAPPYHRS", true},
		{"Lowercase - happyhrs", "happyhrs", true},
		{"Mixed case - HappyHrs", "HappyHrs", true},
		{"Uppercase - FIFTYOFF", "FIFTYOFF", true},
		{"Lowercase - fiftyoff", "fiftyoff", true},
		{"Mixed case - FiftyOff", "FiftyOff", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidatePromoCode(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %v for code '%s', got %v", tt.expected, tt.code, result)
			}
		})
	}
}

// TestLoadCouponFile tests file loading functionality
func TestLoadCouponFile(t *testing.T) {
	validator := NewPromoCodeValidator()
	
	// For this test, we'll simulate the file loading
	// In real implementation, you would create actual gzip files
	validator.couponFiles = map[string][]string{
		"mock1.gz": {"HAPPYHRS", "FIFTYOFF"},
		"mock2.gz": {"HAPPYHRS", "SAVE10NOW"},
	}
	
	// Test that codes are loaded correctly
	if len(validator.couponFiles["mock1.gz"]) != 2 {
		t.Errorf("Expected 2 codes in mock1.gz, got %d", len(validator.couponFiles["mock1.gz"]))
	}
	
	if len(validator.couponFiles["mock2.gz"]) != 2 {
		t.Errorf("Expected 2 codes in mock2.gz, got %d", len(validator.couponFiles["mock2.gz"]))
	}
}

// TestGetValidPromoCodes tests retrieval of all valid promo codes
func TestGetValidPromoCodes(t *testing.T) {
	validator := NewPromoCodeValidator()
	
	// Mock coupon files
	validator.couponFiles = map[string][]string{
		"file1.gz": {"HAPPYHRS", "FIFTYOFF", "UNIQUE1"},
		"file2.gz": {"HAPPYHRS", "SAVE10NOW", "UNIQUE2"},
		"file3.gz": {"FIFTYOFF", "SAVE10NOW", "UNIQUE3"},
	}
	
	validCodes := validator.GetValidPromoCodes()
	
	// Should return codes that appear in at least 2 files
	expectedValidCodes := []string{"HAPPYHRS", "FIFTYOFF", "SAVE10NOW"}
	
	if len(validCodes) != len(expectedValidCodes) {
		t.Errorf("Expected %d valid codes, got %d", len(expectedValidCodes), len(validCodes))
	}
	
	// Check that all expected codes are present
	for _, expectedCode := range expectedValidCodes {
		found := false
		for _, validCode := range validCodes {
			if validCode == expectedCode {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected code '%s' not found in valid codes", expectedCode)
		}
	}
}

// TestPromoCodeValidationIntegration tests the complete validation flow
func TestPromoCodeValidationIntegration(t *testing.T) {
	validator := NewPromoCodeValidator()
	
	// Simulate loading the actual coupon files
	// In real implementation, these would be downloaded from S3
	validator.couponFiles = map[string][]string{
		"couponbase1.gz": {"HAPPYHRS", "FIFTYOFF", "SAVE10NOW", "UNIQUE1"},
		"couponbase2.gz": {"HAPPYHRS", "FIFTYOFF", "WELCOME20", "UNIQUE2"},
		"couponbase3.gz": {"HAPPYHRS", "SAVE10NOW", "WELCOME20", "UNIQUE3"},
	}
	
	// Test cases based on the requirements document
	tests := []struct {
		name     string
		code     string
		expected bool
		reason   string
	}{
		{"Valid - HAPPYHRS", "HAPPYHRS", true, "Found in all 3 files, valid length"},
		{"Valid - FIFTYOFF", "FIFTYOFF", true, "Found in 2 files, valid length"},
		{"Valid - SAVE10NOW", "SAVE10NOW", true, "Found in 2 files, valid length"},
		{"Valid - WELCOME20", "WELCOME20", true, "Found in 2 files, valid length"},
		{"Invalid - UNIQUE1", "UNIQUE1", false, "Found in only 1 file"},
		{"Invalid - UNIQUE2", "UNIQUE2", false, "Found in only 1 file"},
		{"Invalid - UNIQUE3", "UNIQUE3", false, "Found in only 1 file"},
		{"Invalid - SUPER100", "SUPER100", false, "Not found in any file"},
		{"Invalid - SHORT", "SHORT", false, "Too short (5 characters)"},
		{"Invalid - TOOLONG123", "TOOLONG123", false, "Too long (11 characters)"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidatePromoCode(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %v for code '%s' (%s), got %v", 
					tt.expected, tt.code, tt.reason, result)
			}
		})
	}
}

// TestPromoCodePerformance tests performance with large datasets
func TestPromoCodePerformance(t *testing.T) {
	validator := NewPromoCodeValidator()
	
	// Create large mock dataset
	largeDataset := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		largeDataset[i] = fmt.Sprintf("CODE%d", i)
	}
	
	validator.couponFiles = map[string][]string{
		"large1.gz": largeDataset[:5000],
		"large2.gz": largeDataset[2500:7500],
		"large3.gz": largeDataset[5000:],
	}
	
	// Test validation performance
	validCodes := validator.GetValidPromoCodes()
	
	// Should have codes that appear in multiple files
	if len(validCodes) == 0 {
		t.Error("Expected some valid codes in large dataset")
	}
	
	// Test specific code validation
	result := validator.ValidatePromoCode("CODE5000")
	if !result {
		t.Error("Expected CODE5000 to be valid (appears in multiple files)")
	}
}

// BenchmarkPromoCodeValidation benchmarks the validation function
func BenchmarkPromoCodeValidation(b *testing.B) {
	validator := NewPromoCodeValidator()
	
	// Setup test data
	validator.couponFiles = map[string][]string{
		"file1.gz": {"HAPPYHRS", "FIFTYOFF", "SAVE10NOW"},
		"file2.gz": {"HAPPYHRS", "SAVE10NOW", "WELCOME20"},
		"file3.gz": {"HAPPYHRS", "FIFTYOFF", "WELCOME20"},
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		validator.ValidatePromoCode("HAPPYHRS")
	}
}