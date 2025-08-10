package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
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

	// API v1 routes - using a more specific pattern-based approach
	mux.HandleFunc("/api/v1/users/search", h.SearchUsers) // Must come before APIUsersPath
	mux.HandleFunc(APIUsersPath, h.handleUserRoutes)      // Handles all /users/* routes including exact match
	mux.HandleFunc("/api/v1/stats", h.GetStats)
}

// handleUsers handles requests to /api/v1/users (exact match)
func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	h.corsMiddleware(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		h.CreateUser(w, r)
	case "GET":
		h.GetAllUsers(w, r)
	default:
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

// handleUserRoutes handles requests to /api/v1/users/* (with path)
func (h *Handler) handleUserRoutes(w http.ResponseWriter, r *http.Request) {
	h.corsMiddleware(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Extract the path after /api/v1/users/
	path := strings.TrimPrefix(r.URL.Path, APIUsersPath)

	// If path is empty, redirect to exact users route
	if path == "" {
		h.handleUsers(w, r)
		return
	}

	// Check if it's a search request
	if path == "search" {
		h.SearchUsers(w, r)
		return
	}

	// Otherwise, treat it as a user ID request
	switch r.Method {
	case "GET":
		h.GetUser(w, r)
	case "PUT":
		h.UpdateUser(w, r)
	case "DELETE":
		h.DeleteUser(w, r)
	default:
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

// Helper functions for JSON responses and middleware

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// getUserIDFromPath extracts user ID from URL path
func getUserIDFromPath(path string) (int, error) {
	// Remove the /api/v1/users/ prefix
	userIDStr := strings.TrimPrefix(path, APIUsersPath)
	// Remove any trailing slashes or additional path segments
	if idx := strings.Index(userIDStr, "/"); idx != -1 {
		userIDStr = userIDStr[:idx]
	}
	userIDStr = strings.TrimSpace(userIDStr)

	if userIDStr == "" {
		return 0, errors.New("user ID is required")
	}

	// Convert string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, errors.New("user ID must be a valid integer")
	}

	return userID, nil
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

// CreateUser handles user creation requests
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error:   ErrInvalidRequestBody,
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Stub implementation - return not implemented
	writeJSON(w, http.StatusNotImplemented, models.ErrorResponse{
		Error:   ErrNotImplemented,
		Message: ErrServiceNotAvailable,
		Code:    http.StatusNotImplemented,
	})
}

// GetUser handles get user by ID requests
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error:   ErrInvalidUserID,
			Message: ErrUserIDMustBeInt,
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Stub implementation - return not implemented
	writeJSON(w, http.StatusNotImplemented, models.ErrorResponse{
		Error:   ErrNotImplemented,
		Message: ErrServiceNotAvailable,
		Code:    http.StatusNotImplemented,
	})
}

// GetAllUsers handles get all users requests
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Stub implementation - return empty list
	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgUsersRetrieved,
		Data:    []interface{}{},
	})
}

// UpdateUser handles user update requests
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error:   ErrInvalidUserID,
			Message: ErrUserIDMustBeInt,
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error:   ErrInvalidRequestBody,
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Stub implementation - return not implemented
	writeJSON(w, http.StatusNotImplemented, models.ErrorResponse{
		Error:   ErrNotImplemented,
		Message: ErrServiceNotAvailable,
		Code:    http.StatusNotImplemented,
	})
}

// DeleteUser handles user deletion requests
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error:   ErrInvalidUserID,
			Message: ErrUserIDMustBeInt,
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Stub implementation - return not implemented for delete operation
	writeJSON(w, http.StatusNotImplemented, models.ErrorResponse{
		Error:   ErrNotImplemented,
		Message: ErrFailedToDeleteUser + " - " + ErrServiceNotAvailable,
		Code:    http.StatusNotImplemented,
	})
}

// SearchUsers handles user search requests
func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error:   ErrMissingSearchQuery,
			Message: "Query parameter 'q' is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Stub implementation - return empty results
	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgSearchCompleted,
		Data:    []interface{}{},
	})
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
