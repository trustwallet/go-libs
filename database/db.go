package database

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"gorm.io/gorm"
)

//go:generate mockgen -destination=./mock_db.go -package=database . DBContextGetter,TrxContextGetter

type DBContextGetter interface {
	DBFrom(ctx context.Context) *gorm.DB
}

type TrxContextGetter interface {
	Transaction(ctx context.Context, fc func(ctx context.Context) error) error
}

type transactionKey struct{}

var trxKey = &transactionKey{}

// DBGetter implements the DBContextGetter interface, allowing retrieval of a read/write database connection.
type DBGetter struct {
	db *gorm.DB
}

// NewDBGetter creates a new DBGetter instance with the specified database configuration.
// If you are using a read-write splitting database connection, the `dbresolver` will automatically select
// the appropriate connection based on the SQL to be executed.
// When using database transactions, the read-write connection is used by default.
// If you want to force the use of write or read connection, you can use the following method:
// ```
//
//	getter.DBFrom(ctx).Clauses(dbresolver.Write/dbresolver.Read)
//	getter.GetSourceDB().Clauses(dbresolver.Write/dbresolver.Read)
//
// ```
// For more information, read https://gorm.io/docs/dbresolver.html#Read-x2F-Write-Splitting
// *Note* that when using read-write splitting, there is a potential issue where a read operation immediately
// following a write operation may not see the updated data if it is executed on a different read-only replica
// that has not yet been updated with the new data.
func NewDBGetter(cfg DBConfig) (*DBGetter, error) {
	cfg.applyDefaultValue()

	logLevel, err := newLogLevelFromString(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.Url), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormLogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			}),
	})
	if err != nil {
		return nil, err
	}

	var replicas []gorm.Dialector
	if cfg.ReadonlyUrl != nil {
		replicas = append(replicas, postgres.Open(*cfg.ReadonlyUrl))
	}

	resolver := dbresolver.Register(
		dbresolver.Config{
			Sources:           []gorm.Dialector{postgres.Open(cfg.Url)},
			Replicas:          replicas,
			TraceResolverMode: logLevel == gormLogger.Info,
		}).SetConnMaxIdleTime(cfg.ConnPool.ConnMaxIdleTime).
		SetConnMaxLifetime(cfg.ConnPool.ConnMaxLifetime).
		SetMaxIdleConns(cfg.ConnPool.MaxIdleConns).
		SetMaxOpenConns(cfg.ConnPool.MaxOpenConns)
	if err := db.Use(resolver); err != nil {
		return nil, err
	}
	return &DBGetter{db: db}, nil
}

func (getter *DBGetter) GetSourceDB() *gorm.DB {
	return getter.db
}

func (getter *DBGetter) HealthCheck() error {
	// gorm dbresolver doesn't support getting replica connection
	// https://github.com/go-gorm/dbresolver/issues/45
	sqlDB, err := getter.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (getter *DBGetter) DBFrom(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(trxKey).(*gorm.DB); ok {
		return db
	}
	return getter.db
}

func (getter *DBGetter) Transaction(ctx context.Context, fc func(ctx context.Context) error) error {
	return getter.DBFrom(ctx).Transaction(func(tx *gorm.DB) error {
		return fc(context.WithValue(ctx, trxKey, tx))
	})
}

func (getter *DBGetter) Close() error {
	// gorm dbresolver doesn't support getting replica connection
	// https://github.com/go-gorm/dbresolver/issues/45
	sqlDB, err := getter.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
