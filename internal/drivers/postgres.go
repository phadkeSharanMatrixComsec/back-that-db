package drivers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PostgresDriver struct {
	host     string
	port     string
	user     string
	password string
	database string
}

// NewPostgresDriver creates a new PostgreSQL driver instance
func NewPostgresDriver(connString string) (*PostgresDriver, error) {
	// Parse connection string
	// Format: postgres://username:password@host:port/dbname
	// or user=foo password=bar host=baz port=5432 dbname=qux
	var pd PostgresDriver

	if strings.HasPrefix(connString, "postgres://") {
		// URL format
		connString = strings.TrimPrefix(connString, "postgres://")
		parts := strings.Split(connString, "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid connection string format")
		}

		// Parse credentials
		creds := strings.Split(parts[0], ":")
		if len(creds) != 2 {
			return nil, fmt.Errorf("invalid credentials format")
		}
		pd.user = creds[0]
		pd.password = creds[1]

		// Parse host, port, and database
		hostParts := strings.Split(parts[1], "/")
		if len(hostParts) != 2 {
			return nil, fmt.Errorf("invalid host/database format")
		}

		hostPort := strings.Split(hostParts[0], ":")
		pd.host = hostPort[0]
		if len(hostPort) > 1 {
			pd.port = hostPort[1]
		} else {
			pd.port = "5432" // Default PostgreSQL port
		}

		pd.database = hostParts[1]
	} else {
		// Key-value format
		pairs := strings.Split(connString, " ")
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) != 2 {
				continue
			}
			switch kv[0] {
			case "user":
				pd.user = kv[1]
			case "password":
				pd.password = kv[1]
			case "host":
				pd.host = kv[1]
			case "port":
				pd.port = kv[1]
			case "dbname":
				pd.database = kv[1]
			}
		}
	}

	// Validate required fields
	if pd.host == "" || pd.user == "" || pd.database == "" {
		return nil, fmt.Errorf("missing required connection parameters")
	}

	if pd.port == "" {
		pd.port = "5432"
	}

	return &pd, nil
}

// Name returns the name of the driver
func (d *PostgresDriver) Name() string {
	return "postgresql"
}

// Backup performs a database backup using pg_dump
func (d *PostgresDriver) Backup(outPath string) error {
	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build environment variables for authentication
	env := []string{
		fmt.Sprintf("PGUSER=%s", d.user),
		fmt.Sprintf("PGHOST=%s", d.host),
		fmt.Sprintf("PGPORT=%s", d.port),
		fmt.Sprintf("PGDATABASE=%s", d.database),
	}

	if d.password != "" {
		env = append(env, fmt.Sprintf("PGPASSWORD=%s", d.password))
	}

	// Create pg_dump command
	cmd := exec.Command("pg_dump",
		"--format=custom", // Use custom format for better restore control
		"--no-owner",      // Skip restoration of object ownership
		"--no-privileges", // Don't include commands to set privileges
		"--file="+outPath, // Output file
	)

	cmd.Env = append(os.Environ(), env...)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	return nil
}

// Restore performs a database restore using pg_restore
func (d *PostgresDriver) Restore(inPath string) error {
	// Check if input file exists
	if _, err := os.Stat(inPath); err != nil {
		return fmt.Errorf("backup file not found: %w", err)
	}

	// Build environment variables for authentication
	env := []string{
		fmt.Sprintf("PGUSER=%s", d.user),
		fmt.Sprintf("PGHOST=%s", d.host),
		fmt.Sprintf("PGPORT=%s", d.port),
		fmt.Sprintf("PGDATABASE=%s", d.database),
	}

	if d.password != "" {
		env = append(env, fmt.Sprintf("PGPASSWORD=%s", d.password))
	}

	// Create pg_restore command
	cmd := exec.Command("pg_restore",
		"--clean",         // Clean (drop) database objects before recreating
		"--if-exists",     // Add IF EXISTS to DROP commands
		"--no-owner",      // Skip restoration of object ownership
		"--no-privileges", // Don't include commands to set privileges
		"--dbname="+d.database,
		inPath,
	)

	cmd.Env = append(os.Environ(), env...)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_restore failed: %w", err)
	}

	return nil
}
