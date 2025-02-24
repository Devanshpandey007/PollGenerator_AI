// Package S3mock provides a mock implementation of the S3 client for testing.
package S3

import (
	"context"
	"errors"
	"time"
)

// MockS3Client is a mock implementation of S3ClientInterface.
type MockS3Client struct {
	Files        map[string][]byte
	fileListings map[string][]string
	Timeout      time.Duration
}

// NewMockS3Client creates a new instance of MockS3Client.
func NewMockS3Client() *MockS3Client {
	return &MockS3Client{
		Files:        make(map[string][]byte),
		fileListings: make(map[string][]string),
	}
}

// UploadFile mocks uploading a file to S3.
func (m *MockS3Client) UploadFile(ctx context.Context, key string, data []byte) error {
	m.Files[key] = data
	return nil
}

// DownloadFile mocks downloading a file from S3.
func (m *MockS3Client) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	data, exists := m.Files[key]
	if !exists {
		return nil, errors.New("file not found")
	}
	return data, nil
}

// DeleteFile mocks deleting a file from S3.
func (m *MockS3Client) DeleteFile(ctx context.Context, key string) error {
	if _, exists := m.Files[key]; !exists {
		return errors.New("file not found")
	}
	delete(m.Files, key)
	return nil
}

// ListFiles mocks listing files in an S3 bucket.
func (m *MockS3Client) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	var files []string
	for key := range m.Files {
		files = append(files, key)
	}
	return files, nil
}

// Close is a no-op for the mock implementation.
func (m *MockS3Client) Close() error {
	return nil
}

// SetTimeout sets a custom timeout for the mock client.
func (m *MockS3Client) SetTimeout(timeout time.Duration) {
	m.Timeout = timeout
}
