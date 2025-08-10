package api

import (
	"encoding/json"
	"net/http"
	"time"

	"compmodule/internal/models"
)

// Error message constants
const (
	ErrInvalidUserID         = "Invalid user ID"
	ErrUserIDMustBeInt       = "User ID must be a valid integer"
	ErrInvalidRequestBody    = "Invalid request body"
	ErrFailedToCreateUser    = "Failed to create user"
	ErrUserNotFound          = "User not found"
	ErrFailedToRetrieveUsers = "Failed to retrieve users"
	ErrFailedToUpdateUser    = "Failed to update user"
	ErrFailedToDeleteUser    = "Failed to delete user"
	ErrMissingSearchQuery    = "Missing search query"
	ErrSearchFailed          = "Search failed"
	ErrMethodNotAllowed      = "Method not allowed"
	ErrUsernameExists        = "username already exists"
	ErrEmailExists           = "email already exists"
	ErrNotImplemented        = "Not implemented"
	ErrServiceNotAvailable   = "User service not available"
)

// Success message constants
const (
	MsgUserCreated     = "User created successfully"
	MsgUserRetrieved   = "User retrieved successfully"
	MsgUsersRetrieved  = "Users retrieved successfully"
	MsgUserUpdated     = "User updated successfully"
	MsgUserDeleted     = "User deleted successfully"
	MsgSearchCompleted = "Search completed successfully"
	MsgStatsRetrieved  = "Statistics retrieved successfully"
)

// API route constants
const (
	APIUsersPath = "/api/v1/users/"
)

// Handler holds the dependencies for API handlers
type Handler struct {
	// Removed userService dependency for now
}

// NewHandler creates a new Handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// SetupRoutes sets up the API routes
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/health", h.HealthCheck)

	mux.HandleFunc("/api/v1/stats", h.GetStats)
}

// Helper functions for JSON responses and middleware

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// corsMiddleware handles CORS headers
func (h *Handler) corsMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Apply CORS middleware
	h.corsMiddleware(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	health := &models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}
	writeJSON(w, http.StatusOK, health)
}

// GetStats handles statistics requests
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"total_users": 0, // Stub implementation
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgStatsRetrieved,
		Data:    stats,
	})
}

// Note: Middleware functions (CORS and Logging) are now implemented
// as methods on the Handler struct and applied within individual handlers.
// This provides better control and eliminates the dependency on Gin framework.
