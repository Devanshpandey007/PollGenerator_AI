// Package S3 defines a client for interacting with AWS S3.
package S3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/transport/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client provides methods to interact with an S3 bucket.
type S3Client struct {
	Client  *s3.Client
	Timeout time.Duration
	Bucket  string
}

// ErrNoSuchKey is returned when the requested key does not exist in S3.
var ErrNoSuchKey = errors.New("no such key")

// NewS3Client creates a new instance of S3Client.
func NewS3Client(bucket, region string, timeout time.Duration) (*S3Client, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Return the initialized S3Client
	return &S3Client{
		Client:  client,
		Timeout: timeout,
		Bucket:  bucket,
	}, nil
}

// UploadFile uploads a file to the specified S3 bucket.
func (s *S3Client) UploadFile(ctx context.Context, key string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	uploader := manager.NewUploader(s.Client)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return nil
}

// DownloadFile downloads a file from the specified S3 bucket.
func (s *S3Client) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	buffer := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(s.Client)

	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var respErr *http.ResponseError
		if errors.As(err, &respErr) {
			var smithyErr smithy.APIError
			if errors.As(err, &smithyErr) && smithyErr.ErrorCode() == "NoSuchKey" {
				return nil, ErrNoSuchKey
			}
		}
		return nil, fmt.Errorf("failed to download file from S3: %v", err)
	}

	return buffer.Bytes(), nil
}

// DeleteFile deletes a file from the specified S3 bucket.
func (s *S3Client) DeleteFile(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	return nil
}

// ListFiles lists all files in the specified S3 bucket with a given prefix.
func (s *S3Client) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	var files []string
	paginator := s3.NewListObjectsV2Paginator(s.Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list files in S3: %v", err)
		}

		for _, object := range output.Contents {
			files = append(files, *object.Key)
		}
	}

	return files, nil
}

// Close is a placeholder for the S3 client (not needed for S3, but included for consistency).
func (s *S3Client) Close() error {
	// AWS S3 client does not require an explicit close.
	return nil
}

// SetTimeout sets a custom timeout for the S3 client.
func (s *S3Client) SetTimeout(timeout time.Duration) {
	s.Timeout = timeout
}
