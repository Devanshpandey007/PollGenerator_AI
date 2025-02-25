package Api

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock implementation of the HTTPClient interface.
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestMakeGetRequest_Success(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockHTTPClient)

	// Create a mock response
	responseBody := `{"key":"value"}`
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	// Set up expected behavior on mock client
	mockClient.On("Do", mock.Anything).Return(mockResp, nil)
	// Initialize ApiClient with the mock client
	apiClient := &ApiClient{
		Client: mockClient,
		APIKey: "test-api-key",
	}

	// Perform the request
	url := "www.example.com/test-endpoint"
	headers := map[string]string{"Custom-Header": "HeaderValue"}
	body, err := apiClient.MakeGetRequest(ctx, url, headers)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}, body)
	mockClient.AssertExpectations(t)
}

func TestMakeGetRequest_Non200Status(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockHTTPClient)

	// Create a mock response with a non-200 status code
	responseBody := `{"error":"not found"}`
	mockResp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	// Set up expected behavior on mock client
	mockClient.On("Do", mock.Anything).Return(mockResp, nil)

	// Initialize ApiClient with the mock client
	apiClient := &ApiClient{
		Client: mockClient,
		APIKey: "test-api-key",
	}

	// Perform the request
	url := "www.example.com/test-endpoint"
	body, err := apiClient.MakeGetRequest(ctx, url, nil)
	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error: status_code 404")
	assert.Nil(t, body)
	mockClient.AssertExpectations(t)
}

func TestMakeGetRequest_RequestError(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockHTTPClient)

	// Set up the mock to return an error for the request
	mockClient.On("Do", mock.Anything).Return(nil, errors.New("network error"))

	// Initialize ApiClient with the mock client
	apiClient := &ApiClient{
		Client: mockClient,
		APIKey: "test-api-key",
	}

	// Perform the request
	url := "www.example.com/test-endpoint"
	body, err := apiClient.MakeGetRequest(ctx, url, nil)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make request: network error")
	assert.Nil(t, body)
	mockClient.AssertExpectations(t)
}
