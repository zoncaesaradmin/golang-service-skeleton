package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"compmodule/internal/models"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler()
	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}
}

func TestHealthCheck(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.HealthCheck(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("HealthCheck() status = %d, want %d", rr.Code, http.StatusOK)
	}

	// Check Content-Type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("HealthCheck() Content-Type = %q, want %q", contentType, "application/json")
	}

	// Parse and validate JSON response
	var health models.HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode health response: %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("HealthCheck() status = %q, want %q", health.Status, "healthy")
	}

	if health.Version != "1.0.0" {
		t.Errorf("HealthCheck() version = %q, want %q", health.Version, "1.0.0")
	}
}

func TestGetStats(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil)
	rr := httptest.NewRecorder()

	handler.GetStats(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("GetStats() status = %d, want %d", rr.Code, http.StatusOK)
	}

	// Check Content-Type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("GetStats() Content-Type = %q, want %q", contentType, "application/json")
	}

	// Parse and validate JSON response
	var response models.SuccessResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode stats response: %v", err)
	}

	if response.Message != MsgStatsRetrieved {
		t.Errorf("GetStats() message = %q, want %q", response.Message, MsgStatsRetrieved)
	}
}

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	data := models.SuccessResponse{Message: "test", Data: "data"}

	writeJSON(rr, http.StatusOK, data)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("writeJSON() status = %d, want %d", rr.Code, http.StatusOK)
	}

	// Check Content-Type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("writeJSON() Content-Type = %q, want %q", contentType, "application/json")
	}

	// Check that response body is valid JSON
	var result interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Errorf("writeJSON() produced invalid JSON: %v", err)
	}
}

func TestHealthCheckOPTIONS(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodOptions, "/health", nil)
	rr := httptest.NewRecorder()

	handler.HealthCheck(rr, req)

	// Check status code for OPTIONS
	if rr.Code != http.StatusNoContent {
		t.Errorf("HealthCheck OPTIONS status = %d, want %d", rr.Code, http.StatusNoContent)
	}

	// Check CORS headers
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "POST, OPTIONS, GET, PUT, DELETE",
	}

	for header, expectedValue := range expectedHeaders {
		if got := rr.Header().Get(header); got != expectedValue {
			t.Errorf("HealthCheck CORS header %s = %q, want %q", header, got, expectedValue)
		}
	}
}

func TestSetupRoutes(t *testing.T) {
	handler := NewHandler()
	mux := http.NewServeMux()
	
	// Setup routes
	handler.SetupRoutes(mux)
	
	// Test that routes are registered by making requests
	testCases := []struct {
		path           string
		expectedStatus int
	}{
		{"/health", http.StatusOK},
		{"/api/v1/stats", http.StatusOK},
	}
	
	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rr := httptest.NewRecorder()
		
		mux.ServeHTTP(rr, req)
		
		if rr.Code != tc.expectedStatus {
			t.Errorf("Route %s: expected status %d, got %d", tc.path, tc.expectedStatus, rr.Code)
		}
	}
}

func TestGetStatsResponseData(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil)
	rr := httptest.NewRecorder()

	handler.GetStats(rr, req)

	// Parse response
	var response models.SuccessResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode stats response: %v", err)
	}

	// Check that data contains expected fields
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("GetStats() response data is not a map")
	}

	// Verify the data structure
	if totalUsers, exists := data["total_users"]; !exists {
		t.Error("GetStats() response missing 'total_users' field")
	} else if totalUsers != float64(0) {
		t.Errorf("GetStats() total_users = %v, want %v", totalUsers, 0)
	}
}
