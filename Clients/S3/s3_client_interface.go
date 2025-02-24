// Package S3 defines a client interface and implementation for interacting with AWS S3.
package S3

import (
	"context"
	"time"
)

// S3ClientInterface defines the contract for an S3 client.
type S3ClientInterface interface {
	UploadFile(ctx context.Context, key string, data []byte) error
	DownloadFile(ctx context.Context, key string) ([]byte, error)
	DeleteFile(ctx context.Context, key string) error
	ListFiles(ctx context.Context, prefix string) ([]string, error)
	Close() error
	SetTimeout(timeout time.Duration)
}
