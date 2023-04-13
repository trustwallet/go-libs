package database

import (
	"fmt"
	"os"
	"strconv"
)

// Constants representing environment variable keys for migration configuration.
const (
	envKeyPrefix       = "MIGRATION"
	envKeyDsn          = "DSN"
	envKeyOp           = "OPERATION"
	envKeyForceVersion = "FORCE_VERSION"
	envKeyFilesDir     = "FILES_DIR"
)

func envKey(key string) string {
	return fmt.Sprintf("%s_%s", envKeyPrefix, key)
}

func readForceVersion() (int, error) {
	forceVersionRaw, ok := os.LookupEnv(envKey(envKeyForceVersion))
	if !ok {
		return 0, nil
	}

	forceVersion, err := strconv.Atoi(forceVersionRaw)
	if err != nil {
		return 0, fmt.Errorf("convert forceVersion: %v", err)
	}

	return forceVersion, nil
}

// RunMigrationsFromEnv reads migration configuration from environment variables,
// creates a MigrationRunner, and runs the specified migration operation.
func RunMigrationsFromEnv(logger logger) error {
	dsn, ok := os.LookupEnv(envKey(envKeyDsn))
	if !ok {
		return fmt.Errorf("missing env: %s", envKey(envKeyDsn))
	}

	operation, ok := os.LookupEnv(envKey(envKeyOp))
	if !ok {
		return fmt.Errorf("missing env: %s", envKey(envKeyOp))
	}

	forceVersion, err := readForceVersion()
	if err != nil {
		return fmt.Errorf("read forceVersion: %v", err)
	}

	filesDir := os.Getenv(envKey(envKeyFilesDir))

	runner, err := NewMigrationRunner(dsn, WithFilesDir(filesDir), WithLogger(logger))
	if err != nil {
		return fmt.Errorf("new migrations runner: %v", err)
	}

	// Get the current migration version and log it.
	version, dirty, err := runner.Version()
	if err != nil {
		logger.Error(fmt.Sprintf("getting current migration version: %v", err))
	} else {
		logger.Info(fmt.Sprintf("migration version before operation: %d, dirty: %v", version, dirty))
	}

	if err := runner.Run(OperationData{
		ID:           operation,
		ForceVersion: forceVersion,
	}); err != nil {
		return fmt.Errorf("run operation %s: %v", operation, err)
	}

	logger.Info("successfully finished migration")

	// Get the migration version after the operation and log it.
	version, dirty, err = runner.Version()
	if err != nil {
		return fmt.Errorf("getting migration version after operation: %v", err)
	}

	logger.Info(fmt.Sprintf("migration version after operation: %d, dirty: %v", version, dirty))

	return nil
}
