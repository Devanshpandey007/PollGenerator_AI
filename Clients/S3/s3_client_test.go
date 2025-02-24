package S3

import (
	"context"
	"testing"
)

func TestMockS3Client(t *testing.T) {
	mockClient := NewMockS3Client()

	// Test uploading a file
	err := mockClient.UploadFile(context.TODO(), "test-file.txt", []byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to upload file: %v", err)
	}

	// Test downloading the file
	data, err := mockClient.DownloadFile(context.TODO(), "test-file.txt")
	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	if string(data) != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", string(data))
	}

	// Test listing files
	files, err := mockClient.ListFiles(context.TODO(), "")
	if err != nil {
		t.Fatalf("Failed to list files: %v", err)
	}

	if len(files) != 1 || files[0] != "test-file.txt" {
		t.Errorf("Expected file 'test-file.txt', got %v", files)
	}

	// Test deleting the file
	err = mockClient.DeleteFile(context.TODO(), "test-file.txt")
	if err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}

	// Ensure the file is deleted
	_, err = mockClient.DownloadFile(context.TODO(), "test-file.txt")
	if err == nil {
		t.Error("Expected error for downloading deleted file, but got none")
	}
}
