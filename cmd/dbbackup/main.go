package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"back-that-db/internal/backup"
)

func main() {
	var (
		source      string
		target      string
		operation   string
		storageType string
	)

	// Parse command line arguments
	flag.StringVar(&source, "source", "postgres://admin:admin123@localhost:5432/sampledb", "Source database connection string")
	flag.StringVar(&target, "target", "./backups/backup.sql", "Target backup location")
	// db type is inferred from the connection string; no explicit flag
	flag.StringVar(&operation, "op", "restore", "Operation (backup, restore)")
	flag.StringVar(&storageType, "storage", "local", "Storage type (local, s3)")
	flag.Parse()

	if source == "" || target == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Initialize backup service (db type inferred from source)
	service := backup.NewService(source, storageType)

	// Perform operation
	var err error
	switch operation {
	case "backup":
		err = service.Backup(target)
	case "restore":
		err = service.Restore(target)
	default:
		fmt.Printf("Unknown operation: %s\n", operation)
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("Operation failed: %v", err)
	}

	fmt.Println("Operation completed successfully")
}
