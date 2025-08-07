package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a client for the service API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// User represents a user response from the API
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
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

// APIResponse represents a standard API response
type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// HealthCheck performs a health check
func (c *Client) HealthCheck() (*HealthResponse, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return nil, fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode health response: %w", err)
	}

	return &health, nil
}

// CreateUser creates a new user
func (c *Client) CreateUser(req *CreateUserRequest) (*User, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/api/v1/users", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create user request failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create user failed: %s", apiResp.Error)
	}

	// Convert the data interface{} to User
	userData, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	var user User
	if err := json.Unmarshal(userData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetUser retrieves a user by ID
func (c *Client) GetUser(id int) (*User, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v1/users/%d", c.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("get user request failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user failed: %s", apiResp.Error)
	}

	// Convert the data interface{} to User
	userData, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	var user User
	if err := json.Unmarshal(userData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetAllUsers retrieves all users
func (c *Client) GetAllUsers() ([]*User, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/users")
	if err != nil {
		return nil, fmt.Errorf("get all users request failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get all users failed: %s", apiResp.Error)
	}

	// Convert the data interface{} to []*User
	usersData, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal users data: %w", err)
	}

	var users []*User
	if err := json.Unmarshal(usersData, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (c *Client) UpdateUser(id int, req *UpdateUserRequest) (*User, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/users/%d", c.baseURL, id), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("update user request failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update user failed: %s", apiResp.Error)
	}

	// Convert the data interface{} to User
	userData, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	var user User
	if err := json.Unmarshal(userData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// DeleteUser deletes a user by ID
func (c *Client) DeleteUser(id int) error {
	httpReq, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/users/%d", c.baseURL, id), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("delete user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete user failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SearchUsers searches for users
func (c *Client) SearchUsers(query string) ([]*User, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v1/users/search?q=%s", c.baseURL, query))
	if err != nil {
		return nil, fmt.Errorf("search users request failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search users failed: %s", apiResp.Error)
	}

	// Convert the data interface{} to []*User
	usersData, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal users data: %w", err)
	}

	var users []*User
	if err := json.Unmarshal(usersData, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return users, nil
}

// GetStats retrieves service statistics
func (c *Client) GetStats() (map[string]interface{}, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/stats")
	if err != nil {
		return nil, fmt.Errorf("get stats request failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get stats failed: %s", apiResp.Error)
	}

	stats, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected stats data format")
	}

	return stats, nil
}
