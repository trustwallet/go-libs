package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	LogLevelSilent string = "silent"
	LogLevelError  string = "error"
	LogLevelWarn   string = "warn"
	LogLevelInfo   string = "info"
)

func newLogLevelFromString(logLevel string) (logger.LogLevel, error) {
	switch logLevel {
	case LogLevelSilent:
		return logger.Silent, nil
	case LogLevelError:
		return logger.Error, nil
	case LogLevelWarn:
		return logger.Warn, nil
	case LogLevelInfo:
		return logger.Info, nil
	default:
		return 0, fmt.Errorf("invalid log level")
	}
}

// Connect takes two parameters, representing the connection strings for read-write splitting databases.
// In some cases read-write splitting is not necessary and the parameters can be identical.
func Connect(readerDSN, writerDSN string, logLevel string) (*DBGetter, error) {
	writer, err := connect(writerDSN, logLevel)
	if err != nil {
		return nil, err
	}

	if readerDSN == writerDSN {
		return NewDbWrapper(writer, writer), nil
	}

	reader, err := connect(readerDSN, logLevel)
	if err != nil {
		return nil, err
	}
	return NewDbWrapper(reader, writer), nil
}

func connect(dsn string, logLevel string) (*gorm.DB, error) {
	level, err := newLogLevelFromString(logLevel)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  level,
					IgnoreRecordNotFoundError: true,
					Colorful:                  true,
				},
			),
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
