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

// NewDbWrapper creates DBGetter instance.
// Both of reader and writer should not be empty, they can be the same pointer pointing to the same database connection.
// reader == writer, in which case they represent a shared read-write database connection.
// reader != writer, then they represent separate read and write database connections for read-write splitting.
func NewDbWrapper(reader, writer *gorm.DB) *DBGetter {
	return &DBGetter{
		readWriteDb: writer,
		readOnlyDb:  reader,
	}
}

func (wrapper *DBGetter) DBFrom(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(trxKey).(*gorm.DB); ok {
		return db
	}
	return wrapper.readWriteDb
}

func (wrapper *DBGetter) ReadOnlyDB() *gorm.DB {
	return wrapper.readOnlyDb
}

func (wrapper *DBGetter) Transaction(ctx context.Context, fc func(ctx context.Context) error) error {
	return wrapper.DBFrom(ctx).Transaction(func(tx *gorm.DB) error {
		return fc(context.WithValue(ctx, trxKey, tx))
	})
}

func (wrapper *DBGetter) Close() error {
	if wrapper.readOnlyDb != wrapper.readWriteDb {
		readOnlyDb, err := wrapper.readOnlyDb.DB()
		if err != nil {
			return err
		}
		if err := readOnlyDb.Close(); err != nil {
			return err
		}
	}
	readWriteDb, err := wrapper.readWriteDb.DB()
	if err != nil {
		return err
	}
	return readWriteDb.Close()
}
