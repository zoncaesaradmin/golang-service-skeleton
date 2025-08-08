package models

import (
	"testing"
	"time"
)

func TestUserGetID(t *testing.T) {
	user := &User{ID: 123}

	id := user.GetID()

	if id != 123 {
		t.Errorf("expected ID to be 123, got %v", id)
	}
}

func TestUserGetCreatedAt(t *testing.T) {
	now := time.Now()
	user := &User{CreatedAt: now}

	createdAt := user.GetCreatedAt()

	if !createdAt.Equal(now) {
		t.Errorf("expected CreatedAt to be %v, got %v", now, createdAt)
	}
}

func TestUserGetUpdatedAt(t *testing.T) {
	now := time.Now()
	user := &User{UpdatedAt: now}

	updatedAt := user.GetUpdatedAt()

	if !updatedAt.Equal(now) {
		t.Errorf("expected UpdatedAt to be %v, got %v", now, updatedAt)
	}
}

func TestUserSetUpdatedAt(t *testing.T) {
	user := &User{}
	now := time.Now()

	user.SetUpdatedAt(now)

	if !user.UpdatedAt.Equal(now) {
		t.Errorf("expected UpdatedAt to be %v, got %v", now, user.UpdatedAt)
	}
}

func TestCreateUserRequestFields(t *testing.T) {
	const testEmail = "test@example.com"
	req := CreateUserRequest{
		Username:  "testuser",
		Email:     testEmail,
		FirstName: "Test",
		LastName:  "User",
	}

	if req.Username != "testuser" {
		t.Errorf("expected Username to be 'testuser', got %s", req.Username)
	}
	if req.Email != testEmail {
		t.Errorf("expected Email to be '%s', got %s", testEmail, req.Email)
	}
	if req.FirstName != "Test" {
		t.Errorf("expected FirstName to be 'Test', got %s", req.FirstName)
	}
	if req.LastName != "User" {
		t.Errorf("expected LastName to be 'User', got %s", req.LastName)
	}
}

func TestUpdateUserRequestPointerFields(t *testing.T) {
	username := "updateduser"
	email := "updated@example.com"
	firstName := "Updated"
	lastName := "User"

	req := UpdateUserRequest{
		Username:  &username,
		Email:     &email,
		FirstName: &firstName,
		LastName:  &lastName,
	}

	if req.Username == nil || *req.Username != "updateduser" {
		t.Errorf("expected Username to be 'updateduser', got %v", req.Username)
	}
	if req.Email == nil || *req.Email != "updated@example.com" {
		t.Errorf("expected Email to be 'updated@example.com', got %v", req.Email)
	}
	if req.FirstName == nil || *req.FirstName != "Updated" {
		t.Errorf("expected FirstName to be 'Updated', got %v", req.FirstName)
	}
	if req.LastName == nil || *req.LastName != "User" {
		t.Errorf("expected LastName to be 'User', got %v", req.LastName)
	}
}

func TestUpdateUserRequestNilFields(t *testing.T) {
	req := UpdateUserRequest{}

	if req.Username != nil {
		t.Errorf("expected Username to be nil, got %v", req.Username)
	}
	if req.Email != nil {
		t.Errorf("expected Email to be nil, got %v", req.Email)
	}
	if req.FirstName != nil {
		t.Errorf("expected FirstName to be nil, got %v", req.FirstName)
	}
	if req.LastName != nil {
		t.Errorf("expected LastName to be nil, got %v", req.LastName)
	}
}

func TestErrorResponseFields(t *testing.T) {
	err := ErrorResponse{
		Error:   "test error",
		Message: "test message",
		Code:    400,
	}

	if err.Error != "test error" {
		t.Errorf("expected Error to be 'test error', got %s", err.Error)
	}
	if err.Message != "test message" {
		t.Errorf("expected Message to be 'test message', got %s", err.Message)
	}
	if err.Code != 400 {
		t.Errorf("expected Code to be 400, got %d", err.Code)
	}
}

func TestSuccessResponseFields(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := SuccessResponse{
		Message: "success",
		Data:    data,
	}

	if resp.Message != "success" {
		t.Errorf("expected Message to be 'success', got %s", resp.Message)
	}
	if resp.Data == nil {
		t.Error("expected Data to not be nil")
	}
}

func TestHealthResponseFields(t *testing.T) {
	now := time.Now()
	health := HealthResponse{
		Status:    "healthy",
		Timestamp: now,
		Version:   "1.0.0",
	}

	if health.Status != "healthy" {
		t.Errorf("expected Status to be 'healthy', got %s", health.Status)
	}
	if !health.Timestamp.Equal(now) {
		t.Errorf("expected Timestamp to be %v, got %v", now, health.Timestamp)
	}
	if health.Version != "1.0.0" {
		t.Errorf("expected Version to be '1.0.0', got %s", health.Version)
	}
}

func TestUserCompleteStructure(t *testing.T) {
	const testEmail = "test@example.com"
	now := time.Now()
	user := User{
		ID:        1,
		Username:  "testuser",
		Email:     testEmail,
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test all fields are set correctly
	if user.ID != 1 {
		t.Errorf("expected ID to be 1, got %d", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("expected Username to be 'testuser', got %s", user.Username)
	}
	if user.Email != testEmail {
		t.Errorf("expected Email to be '%s', got %s", testEmail, user.Email)
	}
	if user.FirstName != "Test" {
		t.Errorf("expected FirstName to be 'Test', got %s", user.FirstName)
	}
	if user.LastName != "User" {
		t.Errorf("expected LastName to be 'User', got %s", user.LastName)
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt to be %v, got %v", now, user.CreatedAt)
	}
	if !user.UpdatedAt.Equal(now) {
		t.Errorf("expected UpdatedAt to be %v, got %v", now, user.UpdatedAt)
	}
}
