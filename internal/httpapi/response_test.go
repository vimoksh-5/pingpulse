package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseHelper_Success(t *testing.T) {
	helper := NewResponseHelper()
	
	// Create a test response recorder
	w := httptest.NewRecorder()
	
	// Test data
	testData := map[string]interface{}{
		"projects": []string{"Project1", "Project2"},
		"count":    2,
	}
	
	// Call Success method
	helper.Success(w, testData, "Test success message")
	
	// Check HTTP status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	
	// Check Content-Type header
	expectedContentType := "application/json"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}
	
	// Parse response body
	var response StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Check response structure
	if response.StatusCode != 1 {
		t.Errorf("Expected status_code 1, got %d", response.StatusCode)
	}
	
	if response.StatusMessage != "Test success message" {
		t.Errorf("Expected status_message 'Test success message', got '%s'", response.StatusMessage)
	}
	
	// Check data field
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data to be a map, got %T", response.Data)
	}
	
	if count, ok := data["count"].(float64); !ok || int(count) != 2 {
		t.Errorf("Expected count to be 2, got %v", data["count"])
	}
}

func TestResponseHelper_Error(t *testing.T) {
	helper := NewResponseHelper()
	
	// Create a test response recorder
	w := httptest.NewRecorder()
	
	// Call Error method
	helper.Error(w, http.StatusBadRequest, "Test error message", nil)
	
	// Check HTTP status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	// Parse response body
	var response StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Check response structure
	if response.StatusCode != 0 {
		t.Errorf("Expected status_code 0, got %d", response.StatusCode)
	}
	
	if response.StatusMessage != "Test error message" {
		t.Errorf("Expected status_message 'Test error message', got '%s'", response.StatusMessage)
	}
	
	if response.Data != nil {
		t.Errorf("Expected data to be nil, got %v", response.Data)
	}
}

func TestResponseHelper_BadRequest(t *testing.T) {
	helper := NewResponseHelper()
	
	// Create a test response recorder
	w := httptest.NewRecorder()
	
	// Call BadRequest method
	helper.BadRequest(w, "Invalid parameter")
	
	// Check HTTP status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	// Parse response body
	var response StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Check response structure
	if response.StatusCode != 0 {
		t.Errorf("Expected status_code 0, got %d", response.StatusCode)
	}
	
	if response.StatusMessage != "Invalid parameter" {
		t.Errorf("Expected status_message 'Invalid parameter', got '%s'", response.StatusMessage)
	}
}

func TestResponseHelper_Health(t *testing.T) {
	helper := NewResponseHelper()
	
	// Create a test response recorder
	w := httptest.NewRecorder()
	
	// Call Health method
	helper.Health(w)
	
	// Check HTTP status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	
	// Parse response body
	var response StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Check response structure
	if response.StatusCode != 1 {
		t.Errorf("Expected status_code 1, got %d", response.StatusCode)
	}
	
	if response.StatusMessage != "Service is healthy" {
		t.Errorf("Expected status_message 'Service is healthy', got '%s'", response.StatusMessage)
	}
	
	// Check data field
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data to be a map, got %T", response.Data)
	}
	
	if status, ok := data["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected status to be 'healthy', got %v", data["status"])
	}
}
