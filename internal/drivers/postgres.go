package drivers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type PostgresDriver struct {
	host     string
	port     int
	user     string
	password string
	database string
}

// NewPostgresDriver creates a new PostgreSQL driver instance
// NewPostgresDriver creates a new PostgreSQL driver instance from parsed connection info
func NewPostgresDriver(ci *ConnectionInfo) (*PostgresDriver, error) {
	if ci == nil {
		return nil, fmt.Errorf("connection info is nil")
	}

	pd := &PostgresDriver{
		host:     ci.Host,
		port:     ci.Port,
		user:     ci.User,
		password: ci.Password,
		database: ci.Database,
	}

	if pd.host == "" || pd.user == "" || pd.database == "" {
		return nil, fmt.Errorf("missing required connection parameters")
	}

	if pd.port == 0 {
		pd.port = 5432
	}

	return pd, nil
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
		fmt.Sprintf("PGPORT=%d", d.port),
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
		fmt.Sprintf("PGPORT=%d", d.port),
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
