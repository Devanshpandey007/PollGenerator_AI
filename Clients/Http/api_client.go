package Api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient defines the contract for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ApiClient provides methods to interact with external APIs
type ApiClient struct {
	Client  HTTPClient // Interface to decouple from concrete *http.Client
	APIKey  string     // API key for authentication
	Timeout time.Duration
}

// NewApiClient creates a new instance of ApiClient with default timeout.
func NewApiClient(apiKey string) *ApiClient {
	return &ApiClient{
		Client: &http.Client{Timeout: 10 * time.Second}, // Default 10-second timeout
		APIKey: apiKey,
	}
}

// MakeGetRequest sends a GET request to the specified URL and returns the response body as []byte.
func (a *ApiClient) MakeGetRequest(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	// Create a new HTTP request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers (including the API key)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute the HTTP request
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	// Check if response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status_code %d, body: %s", resp.StatusCode, string(body))
	}
	return resp, nil
}

// SetTimeout sets a custom timeout for the HTTP client.
func (a *ApiClient) SetTimeout(timeout time.Duration) {
	a.Client = &http.Client{Timeout: timeout}
}

// SetClient allows setting a custom HTTP client (for testing or other purposes).
func (a *ApiClient) SetClient(client HTTPClient) {
	a.Client = client
}
