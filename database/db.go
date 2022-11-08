package database

import (
	"context"

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

type DBGetter struct {
	globalDb *gorm.DB
}

func NewDbWrapper(db *gorm.DB) *DBGetter {
	return &DBGetter{globalDb: db}
}

func (wrapper *DBGetter) DBFrom(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(trxKey).(*gorm.DB); ok {
		return db
	}
	return wrapper.globalDb
}

func (wrapper *DBGetter) Transaction(ctx context.Context, fc func(ctx context.Context) error) error {
	return wrapper.DBFrom(ctx).Transaction(func(tx *gorm.DB) error {
		return fc(context.WithValue(ctx, trxKey, tx))
	})
}

func (wrapper *DBGetter) Close() error {
	sqlDB, err := wrapper.globalDb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
