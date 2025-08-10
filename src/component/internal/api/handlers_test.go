package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"compmodule/internal/models"
)

const (
	// API endpoints
	healthEndpoint      = "/health"
	usersEndpoint       = "/api/v1/users/"
	statsEndpoint       = "/api/v1/stats"
	searchEndpoint      = "/api/v1/users/search"
	userIDEndpoint      = "/api/v1/users/1"
	userInvalidEndpoint = "/api/v1/users/invalid"

	// HTTP headers
	contentType     = "Content-Type"
	applicationJSON = "application/json"

	// Test messages
	expectedStatusMsg       = "expected status %d, got %d"
	expectedStatusDetailMsg = "expected status %d for %s %s, got %d"
	failedUnmarshalMsg      = "failed to unmarshal response: %v"
	expectedMessageMsg      = "expected message '%s', got '%s'"
	expectedErrorMsg        = "expected error '%s', got '%s'"

	// Test query parameters
	testQuery = "?q=test"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler()

	if handler == nil {
		t.Fatal("expected handler to not be nil")
	}
}

func TestSetupRoutes(t *testing.T) {
	handler := NewHandler()
	mux := http.NewServeMux()

	handler.SetupRoutes(mux)

	// Test that routes are registered by making requests
	testCases := []struct {
		path           string
		method         string
		expectedStatus int
	}{
		{healthEndpoint, "GET", http.StatusOK},
		{usersEndpoint, "GET", http.StatusOK},
		{statsEndpoint, "GET", http.StatusOK},
		{searchEndpoint + testQuery, "GET", http.StatusOK},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, tc.path, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != tc.expectedStatus {
			t.Errorf(expectedStatusDetailMsg, tc.expectedStatus, tc.method, tc.path, rr.Code)
		}
	}
}

func TestHealthCheck(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", healthEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HealthCheck(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(expectedStatusMsg, http.StatusOK, status)
	}

	var health models.HealthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &health)
	if err != nil {
		t.Fatalf(failedUnmarshalMsg, err)
	}

	if health.Status != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", health.Status)
	}
	if health.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", health.Version)
	}
}

func TestHealthCheckOptions(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("OPTIONS", healthEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HealthCheck(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf(expectedStatusMsg, http.StatusNoContent, status)
	}
}

func TestGetAllUsers(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", usersEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(expectedStatusMsg, http.StatusOK, status)
	}

	var response models.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(failedUnmarshalMsg, err)
	}

	if response.Message != MsgUsersRetrieved {
		t.Errorf(expectedMessageMsg, MsgUsersRetrieved, response.Message)
	}

	// Should return empty array for stub implementation
	if response.Data == nil {
		t.Error("expected data to not be nil")
	}
}

func TestCreateUserStub(t *testing.T) {
	handler := NewHandler()

	createReq := models.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
	}

	body, _ := json.Marshal(createReq)
	req, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(contentType, applicationJSON)

	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)

	// Should return not implemented for stub
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf(expectedStatusMsg, http.StatusNotImplemented, status)
	}

	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(failedUnmarshalMsg, err)
	}

	if !strings.Contains(response.Error, "Not implemented") {
		t.Errorf("expected error to contain 'Not implemented', got '%s'", response.Error)
	}
}

func TestCreateUserInvalidJSON(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(contentType, applicationJSON)

	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf(expectedStatusMsg, http.StatusBadRequest, status)
	}

	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(failedUnmarshalMsg, err)
	}

	if response.Error != ErrInvalidRequestBody {
		t.Errorf(expectedErrorMsg, ErrInvalidRequestBody, response.Error)
	}
}

func TestGetUserStub(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", userIDEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUser(rr, req)

	// Should return not implemented for stub
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf(expectedStatusMsg, http.StatusNotImplemented, status)
	}
}

func TestGetUserInvalidID(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", userInvalidEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf(expectedStatusMsg, http.StatusBadRequest, status)
	}

	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(failedUnmarshalMsg, err)
	}

	if response.Error != ErrInvalidUserID {
		t.Errorf(expectedErrorMsg, ErrInvalidUserID, response.Error)
	}
}

func TestUpdateUserStub(t *testing.T) {
	handler := NewHandler()

	updateReq := models.UpdateUserRequest{
		FirstName: stringPtr("Updated"),
	}

	body, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PUT", userIDEndpoint, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(contentType, applicationJSON)

	rr := httptest.NewRecorder()
	handler.UpdateUser(rr, req)

	// Should return not implemented for stub
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf(expectedStatusMsg, http.StatusNotImplemented, status)
	}
}

func TestDeleteUserStub(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("DELETE", userIDEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.DeleteUser(rr, req)

	// Should return not implemented for stub
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf(expectedStatusMsg, http.StatusNotImplemented, status)
	}
}

func TestSearchUsers(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", searchEndpoint+testQuery, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.SearchUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(expectedStatusMsg, http.StatusOK, status)
	}

	var response models.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(failedUnmarshalMsg, err)
	}

	if response.Message != MsgSearchCompleted {
		t.Errorf(expectedMessageMsg, MsgSearchCompleted, response.Message)
	}
}

func TestSearchUsersMissingQuery(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", searchEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.SearchUsers(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf(expectedStatusMsg, http.StatusBadRequest, status)
	}

	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(failedUnmarshalMsg, err)
	}

	if response.Error != ErrMissingSearchQuery {
		t.Errorf(expectedErrorMsg, ErrMissingSearchQuery, response.Error)
	}
}

func TestGetStats(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", statsEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetStats(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(expectedStatusMsg, http.StatusOK, status)
	}

	var response models.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(failedUnmarshalMsg, err)
	}

	if response.Message != MsgStatsRetrieved {
		t.Errorf(expectedMessageMsg, MsgStatsRetrieved, response.Message)
	}

	// Check that data contains total_users
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	if _, exists := data["total_users"]; !exists {
		t.Error("expected data to contain 'total_users'")
	}
}

func TestGetUserIDFromPath(t *testing.T) {
	testCases := []struct {
		path        string
		expectedID  int
		expectError bool
	}{
		{userIDEndpoint, 1, false},
		{"/api/v1/users/123", 123, false},
		{usersEndpoint, 0, true},
		{userInvalidEndpoint, 0, true},
		{"/api/v1/users/1/extra", 1, false}, // Should extract first ID
	}

	for _, tc := range testCases {
		id, err := getUserIDFromPath(tc.path)

		if tc.expectError {
			if err == nil {
				t.Errorf("expected error for path %s, got none", tc.path)
			}
		} else {
			if err != nil {
				t.Errorf("expected no error for path %s, got %v", tc.path, err)
			}
			if id != tc.expectedID {
				t.Errorf("expected ID %d for path %s, got %d", tc.expectedID, tc.path, id)
			}
		}
	}
}

func TestCORSMiddleware(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.corsMiddleware(rr, req)

	// Check CORS headers are set
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "POST, OPTIONS, GET, PUT, DELETE",
	}

	for header, expectedValue := range expectedHeaders {
		if value := rr.Header().Get(header); value != expectedValue {
			t.Errorf("expected header %s to be '%s', got '%s'", header, expectedValue, value)
		}
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
