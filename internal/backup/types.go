package backup

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// BackupConfig contains backup configuration
type BackupConfig struct {
	Compression bool
	Encryption  bool
	MaxSize     int64
}
