package database

import (
	"fmt"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

type LogLevel string

const (
	LogLevelSilent LogLevel = "silent"
	LogLevelError  LogLevel = "error"
	LogLevelWarn   LogLevel = "warn"
	LogLevelInfo   LogLevel = "info"
)

func newLogLevelFromString(logLevel LogLevel) (gormLogger.LogLevel, error) {
	switch logLevel {
	case LogLevelSilent:
		return gormLogger.Silent, nil
	case LogLevelError:
		return gormLogger.Error, nil
	case LogLevelWarn:
		return gormLogger.Warn, nil
	case LogLevelInfo:
		return gormLogger.Info, nil
	default:
		return 0, fmt.Errorf("invalid log level")
	}
}

type DBConnPool struct {
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// DBConfig represents the configuration for a database connection.
type DBConfig struct {
	// Url is the URL of the read-write database instance to connect to.
	Url string `mapstructure:"url"`

	// ReadonlyUrl is the URL of the read-only database instances.
	// This is optional and can be set to nil if read-write splitting is not required.
	ReadonlyUrl *string `mapstructure:"readonly_url"`

	// LogLevel is the logging level for the database connection.
	// Possible values are "silent", "error", "warn", and "info".
	// This is optional and the default value is "error".
	LogLevel LogLevel `mapstructure:"log_level"`

	// ConnPool is the connection pool settings for the database connection.
	// This is optional and can be set to nil if the default connection pool settings are sufficient.
	ConnPool *DBConnPool `mapstructure:"conn_pool"`
}

var (
	defaultMaxIdleConns    = 2
	defaultMaxOpenConns    = 0
	defaultConnMaxIdleTime = time.Duration(0)
	defaultConnMaxLifetime = time.Duration(0)
)

func (cfg *DBConfig) applyDefaultValue() {
	if cfg.LogLevel == "" {
		cfg.LogLevel = LogLevelError
	}
	if cfg.ConnPool == nil {
		// match the default configuration in database/sql
		// https: //github.com/golang/go/blob/198074abd7ec36ee71198a109d98f1ccdb7c5533/src/database/sql/sql.go#L912
		cfg.ConnPool = &DBConnPool{
			MaxIdleConns:    defaultMaxIdleConns,
			ConnMaxIdleTime: defaultConnMaxIdleTime,
			MaxOpenConns:    defaultMaxOpenConns,
			ConnMaxLifetime: defaultConnMaxLifetime,
		}
	}
}
