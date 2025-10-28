package drivers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ConnectionInfo holds parsed connection details common to all drivers
type ConnectionInfo struct {
	DBType   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// ParseConnString parses a connection string in either URL (postgres://...) or
// key=value format (user=... password=... host=... port=... dbname=...).
// It returns a ConnectionInfo with DBType inferred from URL scheme when present.
func ParseConnString(conn string) (*ConnectionInfo, error) {
	ci := &ConnectionInfo{Port: 0}

	// Try URL parsing first
	if strings.Contains(conn, "//") {
		u, err := url.Parse(conn)
		if err == nil && u.Scheme != "" {
			ci.DBType = u.Scheme
			if u.User != nil {
				ci.User = u.User.Username()
				ci.Password, _ = u.User.Password()
			}
			host := u.Host
			// split host:port if present
			if strings.Contains(host, ":") {
				parts := strings.Split(host, ":")
				ci.Host = parts[0]
				if p, err := strconv.Atoi(parts[1]); err == nil {
					ci.Port = p
				}
			} else {
				ci.Host = host
			}
			// path may contain leading /
			db := strings.TrimPrefix(u.Path, "/")
			ci.Database = db
			if ci.Port == 0 {
				// set default port for common DBs
				switch ci.DBType {
				case "postgres", "postgresql":
					ci.Port = 5432
				case "mysql":
					ci.Port = 3306
				case "mssql":
					ci.Port = 1433
				}
			}
			return ci, nil
		}
	}

	// Fallback to key=value parsing
	pairs := strings.Fields(conn)
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		val := kv[1]
		switch key {
		case "user", "username":
			ci.User = val
		case "password":
			ci.Password = val
		case "host":
			ci.Host = val
		case "port":
			if p, err := strconv.Atoi(val); err == nil {
				ci.Port = p
			}
		case "dbname", "database":
			ci.Database = val
		case "type", "dbtype":
			ci.DBType = val
		}
	}

	if ci.Port == 0 && ci.DBType != "" {
		switch ci.DBType {
		case "postgres", "postgresql":
			ci.Port = 5432
		case "mysql":
			ci.Port = 3306
		case "mssql":
			ci.Port = 1433
		}
	}

	// Basic validation
	if ci.Host == "" || ci.User == "" || ci.Database == "" {
		return nil, fmt.Errorf("incomplete connection information")
	}

	return ci, nil
}
