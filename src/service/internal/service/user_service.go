package service

import (
	"errors"
	"sync"
	"time"

	"katharos/service/internal/models"
)

// Error constants
const (
	ErrUserNotFound      = "user not found"
	ErrUsernameExists    = "username already exists"
	ErrEmailExists       = "email already exists"
	ErrUsernameRequired  = "username is required"
	ErrEmailRequired     = "email is required"
	ErrUsernameMinLength = "username must be at least 3 characters long"
	ErrUsernameMaxLength = "username must be less than 50 characters long"
)

// UserService handles user-related business logic
type UserService struct {
	users  map[int]*models.User
	nextID int
	mutex  sync.RWMutex
}

// NewUserService creates a new UserService instance
func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]*models.User),
		nextID: 1,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if username already exists
	for _, user := range s.users {
		if user.Username == req.Username {
			return nil, errors.New(ErrUsernameExists)
		}
		if user.Email == req.Email {
			return nil, errors.New(ErrEmailExists)
		}
	}

	user := &models.User{
		ID:        s.nextID,
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.users[s.nextID] = user
	s.nextID++

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New(ErrUserNotFound)
	}

	return user, nil
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New(ErrUserNotFound)
	}

	// Check for conflicts with other users
	for uid, u := range s.users {
		if uid == id {
			continue
		}
		if req.Username != nil && u.Username == *req.Username {
			return nil, errors.New(ErrUsernameExists)
		}
		if req.Email != nil && u.Email == *req.Email {
			return nil, errors.New(ErrEmailExists)
		}
	}

	// Update fields
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	user.UpdatedAt = time.Now()

	return user, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New(ErrUserNotFound)
	}

	delete(s.users, id)
	return nil
}

// GetUserCount returns the total number of users
func (s *UserService) GetUserCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.users)
}

// SearchUsers searches for users by username or email
func (s *UserService) SearchUsers(query string) ([]*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var results []*models.User
	for _, user := range s.users {
		if user.Username == query || user.Email == query {
			results = append(results, user)
		}
	}

	return results, nil
}

// ValidateUser validates user data
func (s *UserService) ValidateUser(user *models.User) error {
	if user.Username == "" {
		return errors.New(ErrUsernameRequired)
	}
	if user.Email == "" {
		return errors.New(ErrEmailRequired)
	}
	if len(user.Username) < 3 {
		return errors.New(ErrUsernameMinLength)
	}
	if len(user.Username) > 50 {
		return errors.New(ErrUsernameMaxLength)
	}

	return nil
}

// GetHealthStatus returns the health status of the service
func (s *UserService) GetHealthStatus() *models.HealthResponse {
	return &models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}
}
