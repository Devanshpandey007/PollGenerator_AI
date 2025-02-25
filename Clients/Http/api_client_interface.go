package Api

import (
	"context"
	"net/http"
	"time"
)

// ApiClientInterface defines the contract for API-related operations.
type ApiClientInterface interface {
	MakeGetRequest(ctx context.Context, url string, headers map[string]string) (*http.Response, error) // Fixed to return []byte
	MakePostRequest(ctx context.Context, url string, headers map[string]string, body []byte) (*http.Response, error)
	SetTimeout(timeout time.Duration) // Add method to allow setting timeout
	SetClient(client HTTPClient)      // Add method to allow injecting custom HTTP client
}
