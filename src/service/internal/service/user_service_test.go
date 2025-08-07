package service

import (
	"fmt"
	"testing"

	"katharos/service/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser(t *testing.T) {
	service := NewUserService()

	req := &models.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	user, err := service.CreateUser(req)
	require.NoError(t, err)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.FirstName, user.FirstName)
	assert.Equal(t, req.LastName, user.LastName)
	assert.Equal(t, 1, user.ID)
}

func TestUserService_CreateUser_DuplicateUsername(t *testing.T) {
	service := NewUserService()

	req := &models.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	// Create first user
	_, err := service.CreateUser(req)
	require.NoError(t, err)

	// Try to create duplicate user
	req2 := &models.CreateUserRequest{
		Username:  "testuser",
		Email:     "different@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	_, err = service.CreateUser(req2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrUsernameExists)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	service := NewUserService()

	req := &models.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	// Create first user
	_, err := service.CreateUser(req)
	require.NoError(t, err)

	// Try to create user with duplicate email
	req2 := &models.CreateUserRequest{
		Username:  "differentuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	_, err = service.CreateUser(req2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrEmailExists)
}

func TestUserService_GetUser(t *testing.T) {
	service := NewUserService()

	req := &models.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	createdUser, err := service.CreateUser(req)
	require.NoError(t, err)

	user, err := service.GetUser(createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, createdUser.Username, user.Username)
	assert.Equal(t, createdUser.Email, user.Email)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	service := NewUserService()

	_, err := service.GetUser(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrUserNotFound)
}

func TestUserService_UpdateUser(t *testing.T) {
	service := NewUserService()

	req := &models.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	createdUser, err := service.CreateUser(req)
	require.NoError(t, err)

	newEmail := "updated@example.com"
	updateReq := &models.UpdateUserRequest{
		Email: &newEmail,
	}

	updatedUser, err := service.UpdateUser(createdUser.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, newEmail, updatedUser.Email)
	assert.Equal(t, createdUser.Username, updatedUser.Username) // Should not change
}

func TestUserService_DeleteUser(t *testing.T) {
	service := NewUserService()

	req := &models.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	createdUser, err := service.CreateUser(req)
	require.NoError(t, err)

	err = service.DeleteUser(createdUser.ID)
	require.NoError(t, err)

	// Verify user is deleted
	_, err = service.GetUser(createdUser.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrUserNotFound)
}

func TestUserService_GetAllUsers(t *testing.T) {
	service := NewUserService()

	// Create multiple users
	for i := 0; i < 3; i++ {
		req := &models.CreateUserRequest{
			Username:  fmt.Sprintf("user%d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			FirstName: "Test",
			LastName:  "User",
		}
		_, err := service.CreateUser(req)
		require.NoError(t, err)
	}

	users, err := service.GetAllUsers()
	require.NoError(t, err)
	assert.Len(t, users, 3)
}

func TestUserService_SearchUsers(t *testing.T) {
	service := NewUserService()

	req := &models.CreateUserRequest{
		Username:  "searchuser",
		Email:     "search@example.com",
		FirstName: "Search",
		LastName:  "User",
	}

	_, err := service.CreateUser(req)
	require.NoError(t, err)

	// Search by username
	users, err := service.SearchUsers("searchuser")
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "searchuser", users[0].Username)

	// Search by email
	users, err = service.SearchUsers("search@example.com")
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "search@example.com", users[0].Email)
}

func TestUserService_ValidateUser(t *testing.T) {
	service := NewUserService()

	tests := []struct {
		name        string
		user        *models.User
		expectError string
	}{
		{
			name: "valid user",
			user: &models.User{
				Username: "validuser",
				Email:    "valid@example.com",
			},
			expectError: "",
		},
		{
			name: "empty username",
			user: &models.User{
				Username: "",
				Email:    "valid@example.com",
			},
			expectError: ErrUsernameRequired,
		},
		{
			name: "empty email",
			user: &models.User{
				Username: "validuser",
				Email:    "",
			},
			expectError: ErrEmailRequired,
		},
		{
			name: "username too short",
			user: &models.User{
				Username: "ab",
				Email:    "valid@example.com",
			},
			expectError: ErrUsernameMinLength,
		},
		{
			name: "username too long",
			user: &models.User{
				Username: "this_is_a_very_long_username_that_exceeds_fifty_characters",
				Email:    "valid@example.com",
			},
			expectError: ErrUsernameMaxLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateUser(tt.user)
			if tt.expectError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
			}
		})
	}
}

func TestUserService_GetUserCount(t *testing.T) {
	service := NewUserService()

	assert.Equal(t, 0, service.GetUserCount())

	req := &models.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	_, err := service.CreateUser(req)
	require.NoError(t, err)

	assert.Equal(t, 1, service.GetUserCount())
}

func TestUserService_GetHealthStatus(t *testing.T) {
	service := NewUserService()

	health := service.GetHealthStatus()
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "1.0.0", health.Version)
	assert.NotZero(t, health.Timestamp)
}
