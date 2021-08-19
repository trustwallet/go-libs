package logging

import (
	"github.com/heirko/go-contrib/logrusHelper"
	mate "github.com/heralight/logrus_mate"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

const FieldKeyComponent = "component"

type Config mate.LoggerConfig

func init() {
	logger = logrus.New()
}

func SetLoggerConfig(config Config) error {
	err := logrusHelper.SetConfig(logrus.StandardLogger(), mate.LoggerConfig(config))
	if err != nil {
		return err
	}
	return logrusHelper.SetConfig(logger, mate.LoggerConfig(config))
}

// GetLogger returns the logger instance.
func GetLogger() *logrus.Logger {
	return logger
}

// GetLoggerForComponent returns the logger instance with component field set
func GetLoggerForComponent(component string) *logrus.Entry {
	return GetLogger().WithField(FieldKeyComponent, component)
}

// SetLogger sets the logger instance
// This is useful in testing as the logger can be overridden
// with a test logger
func SetLogger(l *logrus.Logger) {
	logger = l
}
