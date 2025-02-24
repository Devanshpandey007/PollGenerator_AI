package Api

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockApiClient is a mock implementation of the ApiClientInterface.
type MockApiClient struct {
	mock.Mock
}

// MakeGetRequest is a mock implementation of the MakeGetRequest method.
func (m *MockApiClient) MakeGetRequest(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	args := m.Called(ctx, url, headers)
	// Safely get the []byte response, return nil if not set
	var response *http.Response
	if args.Get(0) != nil {
		response = args.Get(0).(*http.Response)
	}
	return response, args.Error(1)
}

// SetTimeout is a mock implementation of the SetTimeout method.
func (m *MockApiClient) SetTimeout(timeout time.Duration) {
	m.Called(timeout)
}

// SetClient is a mock implementation of the SetClient method.
func (m *MockApiClient) SetClient(client HTTPClient) {
	m.Called(client)
}

// FetchEvents is a mock implementation of the FetchEvents method.
func (m *MockApiClient) FetchEvents(ctx context.Context, sportKey string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, sportKey)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// FetchOdds is a mock implementation of the FetchOdds method.
func (m *MockApiClient) FetchOdds(ctx context.Context, sportKey, eventId, bookMaker string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, sportKey, eventId, bookMaker)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// Ensure MockApiClient implements ApiClientInterface at compile time.
var _ ApiClientInterface = (*MockApiClient)(nil)
