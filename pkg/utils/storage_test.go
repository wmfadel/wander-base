package utils

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "valid base URL",
			baseURL: "http://localhost:8080",
		},
		{
			name:    "empty base URL",
			baseURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewStorage(tt.baseURL)
			if storage == nil {
				t.Error("NewStorage() returned nil")
			}
			if storage.baseURL != tt.baseURL {
				t.Errorf("NewStorage() baseURL = %v, want %v", storage.baseURL, tt.baseURL)
			}
		})
	}
}

func TestUploadFile(t *testing.T) {
	// Create a test storage instance
	storage := NewStorage("http://localhost:8080")

	// Create test directories
	testDirs := []string{"public", "public/user_photos", "public/events"}
	for _, dir := range testDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}
	defer os.RemoveAll("public")

	// Create a test file
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	// Create a multipart file header
	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	fileHeader := &multipart.FileHeader{
		Filename: testFile,
		Header:   make(map[string][]string),
	}
	fileHeader.Header.Set("Content-Type", "image/jpeg")

	tests := []struct {
		name    string
		prefix  string
		id      int64
		wantErr bool
	}{
		{
			name:    "valid upload",
			prefix:  "user_photos",
			id:      1,
			wantErr: false,
		},
		{
			name:    "invalid prefix",
			prefix:  "invalid_prefix",
			id:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := storage.UploadFile(fileHeader, tt.prefix, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && url == "" {
				t.Error("UploadFile() returned empty URL")
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	// Create a test storage instance
	storage := NewStorage("http://localhost:8080")

	// Create test file
	testFile := "public/test.txt"
	if err := os.MkdirAll(filepath.Dir(testFile), 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.RemoveAll("public")

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid file",
			url:     "http://localhost:8080/test.txt",
			wantErr: false,
		},
		{
			name:    "non-existent file",
			url:     "http://localhost:8080/nonexistent.txt",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			url:     "invalid-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.DeleteFile(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
