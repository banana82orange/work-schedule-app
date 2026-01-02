package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements FileStorage for local filesystem
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new LocalStorage
func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	// Ensure directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

// Save saves a file to local storage
func (s *LocalStorage) Save(ctx context.Context, fileName string, data []byte) (string, error) {
	filePath := filepath.Join(s.basePath, fileName)

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write data
	if _, err := file.Write(data); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return URL
	return s.baseURL + "/" + fileName, nil
}

// Delete deletes a file from local storage
func (s *LocalStorage) Delete(ctx context.Context, fileURL string) error {
	// Extract filename from URL
	fileName := filepath.Base(fileURL)
	filePath := filepath.Join(s.basePath, fileName)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, consider it deleted
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Get retrieves a file from local storage
func (s *LocalStorage) Get(ctx context.Context, fileURL string) ([]byte, error) {
	// Extract filename from URL
	fileName := filepath.Base(fileURL)
	filePath := filepath.Join(s.basePath, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}
