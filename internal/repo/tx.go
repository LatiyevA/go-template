package repo

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type (
	ctxKey string
)

var transaction ctxKey = "transaction"

func (r *Repository) Begin(ctx context.Context) context.Context {
	db := r.db.Begin()
	return setDBToCTX(ctx, db)
}

func (r *Repository) Commit(ctx context.Context) error {
	db := getDBFromCTX(ctx, r.db)
	return db.Commit().Error
}

func (r *Repository) Rollback(ctx context.Context) error {
	db := getDBFromCTX(ctx, r.db)
	return db.Rollback().Error
}

func setDBToCTX(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, transaction, db)
}

func getDBFromCTX(ctx context.Context, db *gorm.DB) *gorm.DB {
	ctxValueAny := ctx.Value(transaction)
	if ctxValueAny == nil {
		log.Println("no transaction found in context")
		return db
	}

	val, ok := ctxValueAny.(*gorm.DB)
	if !ok {
		log.Println("no transaction found in context")
		return db
	}

	return val
}
