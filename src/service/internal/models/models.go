package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username" binding:"required"`
	Email     string    `json:"email" db:"email" binding:"required,email"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GetID returns the user ID (implements Resource interface)
func (u *User) GetID() interface{} {
	return u.ID
}

// GetCreatedAt returns the creation time (implements Resource interface)
func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// GetUpdatedAt returns the last update time (implements Resource interface)
func (u *User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

// SetUpdatedAt sets the last update time (implements Resource interface)
func (u *User) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = t
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}
