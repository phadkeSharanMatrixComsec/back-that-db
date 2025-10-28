# Back That DB

A command-line tool for backing up and restoring databases to various storage backends.

## Features

- Supports multiple database types:
  - MySQL
  - PostgreSQL
  - Microsoft SQL Server

- Storage backends:
  - Local filesystem
  - S3 (Amazon Simple Storage Service)

## Installation

```bash
go install github.com/yourusername/back-that-db/cmd/dbbackup@latest
```

## Usage

### Backup a database

```bash
dbbackup -source "user:pass@tcp(localhost:3306)/mydb" -target "/path/to/backup.sql" -type mysql -op backup
```

### Restore a database

```bash
dbbackup -source "/path/to/backup.sql" -target "user:pass@tcp(localhost:3306)/mydb" -type mysql -op restore
```

### Using S3 storage

```bash
dbbackup -source "user:pass@tcp(localhost:3306)/mydb" -target "s3://mybucket/backup.sql" -type mysql -storage s3 -op backup
```

## Configuration

### Environment Variables

- `AWS_ACCESS_KEY_ID` - AWS access key for S3 storage
- `AWS_SECRET_ACCESS_KEY` - AWS secret key for S3 storage
- `AWS_REGION` - AWS region for S3 storage
- `AWS_BUCKET` - Default S3 bucket name

## Requirements

- Go 1.19 or later
- Database CLI tools installed for the databases you want to work with:
  - `mysqldump` and `mysql` for MySQL
  - `pg_dump` and `psql` for PostgreSQL
  - `sqlcmd` for Microsoft SQL Server

## Building from source

```bash
git clone https://github.com/yourusername/back-that-db.git
cd back-that-db
go build ./cmd/dbbackup
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

MIT