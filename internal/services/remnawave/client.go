package remnawave

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client представляет клиент для работы с Remnawave API
type Client struct {
	baseURL    string
	apiKey     string
	secretKey  string
	httpClient *http.Client
}

// NewClient создает новый клиент Remnawave
func NewClient(baseURL, apiKey, secretKey string) *Client {
	return &Client{
		baseURL:   baseURL,
		apiKey:    apiKey,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// APIResponse представляет базовый ответ API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Server представляет сервер в Remnawave
type Server struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// Plan представляет тарифный план в Remnawave
type Plan struct {
	ID          int     `json:"id"`
	ServerID    int     `json:"server_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`
	IsActive    bool    `json:"is_active"`
}

// User представляет пользователя в Remnawave
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}

// Subscription представляет подписку в Remnawave
type Subscription struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	ServerID   int       `json:"server_id"`
	PlanID     int       `json:"plan_id"`
	Status     string    `json:"status"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetServers получает список серверов
func (c *Client) GetServers() ([]Server, error) {
	var response struct {
		APIResponse
		Data []Server `json:"data"`
	}

	if err := c.makeRequest("GET", "/servers", nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get servers: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return response.Data, nil
}

// GetPlans получает список тарифных планов
func (c *Client) GetPlans(serverID int) ([]Plan, error) {
	var response struct {
		APIResponse
		Data []Plan `json:"data"`
	}

	url := fmt.Sprintf("/servers/%d/plans", serverID)
	if err := c.makeRequest("GET", url, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return response.Data, nil
}

// CreateUser создает пользователя
func (c *Client) CreateUser(username, email string) (*User, error) {
	var response struct {
		APIResponse
		Data User `json:"data"`
	}

	requestData := map[string]string{
		"username": username,
		"email":    email,
	}

	if err := c.makeRequest("POST", "/users", requestData, &response); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return &response.Data, nil
}

// CreateSubscription создает подписку
func (c *Client) CreateSubscription(userID, serverID, planID int) (*Subscription, error) {
	var response struct {
		APIResponse
		Data Subscription `json:"data"`
	}

	requestData := map[string]int{
		"user_id":   userID,
		"server_id": serverID,
		"plan_id":   planID,
	}

	if err := c.makeRequest("POST", "/subscriptions", requestData, &response); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return &response.Data, nil
}

// GetSubscription получает подписку по ID
func (c *Client) GetSubscription(subscriptionID int) (*Subscription, error) {
	var response struct {
		APIResponse
		Data Subscription `json:"data"`
	}

	url := fmt.Sprintf("/subscriptions/%d", subscriptionID)
	if err := c.makeRequest("GET", url, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return &response.Data, nil
}

// UpdateSubscription обновляет подписку
func (c *Client) UpdateSubscription(subscriptionID int, data map[string]interface{}) (*Subscription, error) {
	var response struct {
		APIResponse
		Data Subscription `json:"data"`
	}

	url := fmt.Sprintf("/subscriptions/%d", subscriptionID)
	if err := c.makeRequest("PUT", url, data, &response); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return &response.Data, nil
}

// DeleteSubscription удаляет подписку
func (c *Client) DeleteSubscription(subscriptionID int) error {
	var response APIResponse

	url := fmt.Sprintf("/subscriptions/%d", subscriptionID)
	if err := c.makeRequest("DELETE", url, nil, &response); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// makeRequest выполняет HTTP запрос к API
func (c *Client) makeRequest(method, endpoint string, data interface{}, result interface{}) error {
	var body io.Reader

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal request data: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("X-API-Key", c.apiKey)

	// Добавляем secret key если он есть
	if c.secretKey != "" {
		req.Header.Set("X-Secret-Key", c.secretKey)
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Проверяем статус код
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Парсим JSON ответ
	if err := json.Unmarshal(responseBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
