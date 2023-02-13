package database

import (
	"context"

	"gorm.io/gorm"
)

//go:generate mockgen -destination=./mock_db.go -package=database . DBContextGetter,TrxContextGetter

type DBContextGetter interface {
	DBFrom(ctx context.Context) *gorm.DB
	ReadOnlyDB() *gorm.DB
}

type TrxContextGetter interface {
	Transaction(ctx context.Context, fc func(ctx context.Context) error) error
}

type transactionKey struct{}

var trxKey = &transactionKey{}

// DBGetter implements the DBContextGetter interface, allowing retrieval of a read/write database connection.
type DBGetter struct {
	readWriteDb *gorm.DB
	readOnlyDb  *gorm.DB
}

// NewDBGetter creates DBGetter instance.
// Both of reader and writer should not be empty, they can be the same pointer pointing to the same database connection.
// reader == writer, in which case they represent a shared read-write database connection.
// reader != writer, then they represent separate read and write database connections for read-write splitting.
func NewDBGetter(reader, writer *gorm.DB) *DBGetter {
	return &DBGetter{
		readWriteDb: writer,
		readOnlyDb:  reader,
	}
}

func (getter *DBGetter) HealthCheck() error {
	for _, db := range getter.allDBs() {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Ping(); err != nil {
			return err
		}
	}
	return nil
}

func (getter *DBGetter) DBFrom(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(trxKey).(*gorm.DB); ok {
		return db
	}
	return getter.readWriteDb
}

func (getter *DBGetter) ReadOnlyDB() *gorm.DB {
	return getter.readOnlyDb
}

func (getter *DBGetter) Transaction(ctx context.Context, fc func(ctx context.Context) error) error {
	return getter.DBFrom(ctx).Transaction(func(tx *gorm.DB) error {
		return fc(context.WithValue(ctx, trxKey, tx))
	})
}

func (getter *DBGetter) Close() error {
	for _, db := range getter.allDBs() {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (getter *DBGetter) allDBs() []*gorm.DB {
	result := []*gorm.DB{getter.readWriteDb}
	if getter.readOnlyDb != getter.readWriteDb {
		result = append(result, getter.readOnlyDb)
	}
	return result
}
