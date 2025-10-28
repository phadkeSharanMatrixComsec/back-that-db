package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type StorageBackend interface {
	Store(sourcePath, targetPath string) error
	Retrieve(sourcePath string) (string, error)
}

type LocalStorage struct{}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (s *LocalStorage) Store(sourcePath, targetPath string) error {
	// Create target directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Copy file
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	target, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Retrieve(sourcePath string) (string, error) {
	// For local storage, we can just return the path as is
	if _, err := os.Stat(sourcePath); err != nil {
		return "", fmt.Errorf("backup file not found: %w", err)
	}
	return sourcePath, nil
}
