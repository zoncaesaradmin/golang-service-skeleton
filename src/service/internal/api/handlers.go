package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"katharos/service/internal/models"
	"katharos/service/internal/service"
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

// Handler holds the dependencies for API handlers
type Handler struct {
	userService *service.UserService
}

// NewHandler creates a new Handler instance
func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// SetupRoutes sets up the API routes
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/health", h.HealthCheck)

	// API v1 routes - using a more specific pattern-based approach
	mux.HandleFunc("/api/v1/users/search", h.SearchUsers) // Must come before /api/v1/users/
	mux.HandleFunc("/api/v1/users/", h.handleUserRoutes)  // Handles all /users/* routes
	mux.HandleFunc("/api/v1/users", h.handleUsers)        // Handles exact /users route
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
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")

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
	userIDStr := strings.TrimPrefix(path, "/api/v1/users/")
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

	health := h.userService.GetHealthStatus()
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

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		writeJSON(w, http.StatusConflict, models.ErrorResponse{
			Error:   ErrFailedToCreateUser,
			Message: err.Error(),
			Code:    http.StatusConflict,
		})
		return
	}

	writeJSON(w, http.StatusCreated, models.SuccessResponse{
		Message: MsgUserCreated,
		Data:    user,
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

	user, err := h.userService.GetUser(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{
			Error:   ErrUserNotFound,
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgUserRetrieved,
		Data:    user,
	})
}

// GetAllUsers handles get all users requests
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{
			Error:   ErrFailedToRetrieveUsers,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgUsersRetrieved,
		Data:    users,
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

	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == service.ErrUserNotFound {
			status = http.StatusNotFound
		} else if err.Error() == service.ErrUsernameExists || err.Error() == service.ErrEmailExists {
			status = http.StatusConflict
		}

		writeJSON(w, status, models.ErrorResponse{
			Error:   ErrFailedToUpdateUser,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgUserUpdated,
		Data:    user,
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

	err = h.userService.DeleteUser(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == service.ErrUserNotFound {
			status = http.StatusNotFound
		}

		writeJSON(w, status, models.ErrorResponse{
			Error:   ErrFailedToDeleteUser,
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgUserDeleted,
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

	users, err := h.userService.SearchUsers(query)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{
			Error:   ErrSearchFailed,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgSearchCompleted,
		Data:    users,
	})
}

// GetStats handles statistics requests
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"total_users": h.userService.GetUserCount(),
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse{
		Message: MsgStatsRetrieved,
		Data:    stats,
	})
}

// Note: Middleware functions (CORS and Logging) are now implemented
// as methods on the Handler struct and applied within individual handlers.
// This provides better control and eliminates the dependency on Gin framework.
