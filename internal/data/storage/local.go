package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements Storage interface for local filesystem.
type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) fullPath(key string) string {
	return filepath.Join(s.basePath, key)
}

// Put writes data from reader to the specified key.
func (s *LocalStorage) Put(ctx context.Context, key string, r io.Reader) error {
	fullPath := s.fullPath(key)

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create and write file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Get returns a reader for the specified key.
func (s *LocalStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := s.fullPath(key)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

// Delete removes the file at the specified key.
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := s.fullPath(key)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// Exists checks if a file exists at the specified key.
func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := s.fullPath(key)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check existence: %w", err)
}
