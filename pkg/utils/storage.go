package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Storage struct {
	baseURL string // e.g., "http://localhost:8080"
}

func NewStorage(baseURL string) *Storage {
	return &Storage{baseURL: baseURL}
}

func (s *Storage) UploadFile(file *multipart.FileHeader, prefix string, id int64) (string, error) {
	// Validate MIME type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return "", fmt.Errorf("only image files are allowed")
	}

	// Ensure directory exists
	dir := filepath.Join("public", prefix)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Generate unique filename
	timestamp := time.Now().Format("20060102T150405")
	filename := fmt.Sprintf("%d-%s", id, timestamp)
	path := filepath.Join(dir, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", path, err)
	}

	// Construct URL
	url := fmt.Sprintf("%s/%s/%s", s.baseURL, prefix, filename)
	return url, nil
}

func (s *Storage) DeleteFile(url string) error {
	// Extract path from URL
	path := strings.TrimPrefix(url, s.baseURL)
	if path == url { // No prefix removed
		return fmt.Errorf("invalid URL format: %s", url)
	}
	fullPath := filepath.Join("public", strings.TrimPrefix(path, "/"))
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", fullPath, err)
	}
	return nil
}
