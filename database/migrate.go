package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	// required by migrate.New... to parse migration files directory
	_ "github.com/golang-migrate/migrate/v4/source/file"

	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
)

// supported migrate operations
const (
	defaultFilesDir = "dbmigrations"
	operationUp     = "up"
	operationDown   = "down"
	operationForce  = "force"
)

type operationFn func(*MigrationRunner, OperationData) error

var supportedOperations = map[string]operationFn{
	operationUp:    runUp,
	operationDown:  runDown,
	operationForce: runForce,
}

// OperationData contains information about the migration operation to be performed.
type OperationData struct {
	ID           string
	ForceVersion int
}

// MigrationRunner is responsible for managing and running database migrations.
type MigrationRunner struct {
	mgr      *migrate.Migrate
	filesDir string
	logger   logger
}

// Option represents a function that configures a MigrationRunner.
type Option func(runner *MigrationRunner)

// WithLogger sets a custom logger for the MigrationRunner.
func WithLogger(logger logger) Option {
	return func(runner *MigrationRunner) {
		runner.logger = logger
	}
}

// WithFilesDir sets a custom directory containing migration files for the MigrationRunner.
func WithFilesDir(filesDir string) Option {
	return func(runner *MigrationRunner) {
		runner.filesDir = filesDir
	}
}

// NewMigrationRunner creates a new MigrationRunner with the given database connection string (dsn) and options.
func NewMigrationRunner(dsn string, opts ...Option) (*MigrationRunner, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open dsn: %v", err)
	}

	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return nil, fmt.Errorf("building postgres driver for db migrations: %w", err)
	}

	runner := &MigrationRunner{}

	for _, opt := range opts {
		opt(runner)
	}

	if runner.filesDir == "" {
		runner.filesDir = defaultFilesDir
	}

	if runner.logger == nil {
		runner.logger = &noopLogger{}
	}

	mgr, err := migrate.NewWithDatabaseInstance(
		"file://"+runner.filesDir,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("creating Migrate object: %w", err)
	}

	mgr.Log = toMigrationsLogger(runner.logger)
	runner.mgr = mgr

	return runner, nil
}

// Run executes the migration operation specified by the OperationData.
func (m *MigrationRunner) Run(operation OperationData) error {
	operationName := operation.ID
	operationFn, found := supportedOperations[operationName]
	if !found {
		return fmt.Errorf("unsupported migration operation: %s", operationName)
	}

	if err := operationFn(m, operation); err != nil {
		return fmt.Errorf("operation %s failed: %v", operationName, err)
	}

	return nil
}

// Version returns the current migration version, a dirty flag, and an error if any.
func (m *MigrationRunner) Version() (version uint, dirty bool, err error) {
	return m.mgr.Version()
}

// runUp runs the "up" migration operation, applying new migrations to the database.
func runUp(m *MigrationRunner, _ OperationData) error {
	m.logger.Info("running migrate UP")

	err := m.mgr.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		m.logger.Info(fmt.Sprintf("no new migrations found in: %s", m.filesDir))
		return nil
	}
	if err != nil {
		return fmt.Errorf("running migrations UP failed: %w", err)
	}
	return nil
}

// runDown runs the "down" migration operation, rolling back the latest applied migration.
func runDown(m *MigrationRunner, _ OperationData) error {
	m.logger.Info(fmt.Sprintf("running migrate DOWN with STEPS=%d", 1))

	// always rollback the latest applied migration only
	err := m.mgr.Steps(-1)
	if err != nil {
		return fmt.Errorf("running migrations DOWN failed: %w", err)
	}
	return nil
}

// runForce runs the "force" migration operation, forcibly setting the migration version without running the actual migrations.
func runForce(m *MigrationRunner, op OperationData) error {
	m.logger.Info(fmt.Sprintf("running FORCE with VERSION %d", op.ForceVersion))

	err := m.mgr.Force(op.ForceVersion)
	if err != nil {
		return fmt.Errorf("running migrations FORCE with VERSION %d failed: %w", op.ForceVersion, err)
	}
	return nil
}

type logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Printf(format string, v ...interface{})
}

type noopLogger struct{}

func (l *noopLogger) Info(...interface{})           {}
func (l *noopLogger) Error(...interface{})          {}
func (l *noopLogger) Printf(string, ...interface{}) {}
func (l *noopLogger) Verbose() bool                 { return false }

// adapter to use logger like logrus
func toMigrationsLogger(logger logger) *migrationsLogger {
	return &migrationsLogger{logger: logger}
}

// to be able to log not only errors, but also Info level logs from golang-migrate,
// we have to implement migrate.Logger interface
type migrationsLogger struct {
	logger logger
}

// Printf is like fmt.Printf
func (m *migrationsLogger) Printf(format string, v ...interface{}) {
	m.logger.Info(format, v)
}

// Verbose should return true when verbose logging output is wanted
func (m *migrationsLogger) Verbose() bool {
	return true
}

func (m *migrationsLogger) Info(args ...interface{}) {
	m.logger.Info(args)
}

func (m *migrationsLogger) Error(args ...interface{}) {
	m.logger.Error(args)
}
