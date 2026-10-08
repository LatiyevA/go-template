package repo

import (
	"context"

	"gorm.io/gorm"
)

// WithTransaction commits a successful callback and rolls back on an error or
// panic. The callback's repository is bound to the transaction; nested calls
// use GORM savepoints. The context controls transaction cancellation.
func (r *Repository) WithTransaction(ctx context.Context, fn func(*Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx, log: r.log})
	})
}
