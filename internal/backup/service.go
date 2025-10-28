package backup

import (
	"fmt"
	"os"

	"back-that-db/internal/drivers"
	"back-that-db/internal/storage"
)

type Service struct {
	driver  drivers.DatabaseDriver
	storage storage.StorageBackend
}

// NewService creates a new backup service with the specified database and storage type
func NewService(dbType, storageType, connectionString string) *Service {
	var driver drivers.DatabaseDriver
	var backend storage.StorageBackend

	// Initialize database driver
	switch dbType {
	// case "mysql":
	// 	driver = drivers.NewMySQLDriver(connectionString)
	case "postgres":
		driver, _ = drivers.NewPostgresDriver(connectionString)
	// case "mssql":
	// 	driver, _ = drivers.NewMSSQLDriver(connectionString)
	default:
		panic(fmt.Sprintf("unsupported database type: %s", dbType))
	}

	// Initialize storage backend
	switch storageType {
	case "local":
		backend = storage.NewLocalStorage()
	case "s3":
		backend = storage.NewS3Storage()
	default:
		panic(fmt.Sprintf("unsupported storage type: %s", storageType))
	}

	return &Service{
		driver:  driver,
		storage: backend,
	}
}

// Backup performs a database backup
func (s *Service) Backup(target string) error {
	// Create temporary file for backup
	tmpFile, err := os.CreateTemp("", "db-backup-*.dump")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temp file when done
	defer tmpFile.Close()

	err = s.driver.Backup(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	// Store the backup file
	err = s.storage.Store(tmpFile.Name(), target)
	if err != nil {
		return fmt.Errorf("storing backup failed: %w", err)
	}

	return nil
}

// Restore performs a database restore
func (s *Service) Restore(target string) error {
	// Retrieve backup file
	backupFile, err := s.storage.Retrieve(target)
	if err != nil {
		return fmt.Errorf("retrieving backup failed: %w", err)
	}

	// Restore from backup file
	err = s.driver.Restore(backupFile)
	if err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	return nil
}
