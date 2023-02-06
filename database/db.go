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

type DBGetter struct {
	readWriteDb *gorm.DB
	readOnlyDb  *gorm.DB
}

func NewDbWrapper(db *gorm.DB) *DBGetter {
	return &DBGetter{readWriteDb: db}
}

func NewReadWriteDbWrapper(reader, writer *gorm.DB) *DBGetter {
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
	if wrapper.readOnlyDb != nil {
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
