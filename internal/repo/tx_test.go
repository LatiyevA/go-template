package repo

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/LatiyevA/go-template/pkg/gorm/postgres"
)

func transactionRepository(t *testing.T) *Repository {
	t.Helper()
	url := os.Getenv("TEST_PG_URL")
	if url == "" {
		t.Skip("set TEST_PG_URL to run PostgreSQL integration tests")
	}
	db, err := postgres.New(context.Background(), url, postgres.MaxOpenConns(1), postgres.MaxIdleConns(1), postgres.SilentLogger())
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Error(err)
		}
	})
	if err := db.Exec("CREATE TEMP TABLE tx_items (id bigint PRIMARY KEY)").Error; err != nil {
		t.Fatal(err)
	}
	return New(db, slog.New(slog.NewJSONHandler(io.Discard, nil)))
}

func itemCount(t *testing.T, r *Repository) int64 {
	t.Helper()
	var count int64
	if err := r.db.Raw("SELECT count(*) FROM tx_items").Scan(&count).Error; err != nil {
		t.Fatal(err)
	}
	return count
}

func TestTransactionCommitAndRollback(t *testing.T) {
	r := transactionRepository(t)
	ctx := context.Background()
	if err := r.WithTransaction(ctx, func(tx *Repository) error {
		return tx.db.Exec("INSERT INTO tx_items VALUES (1)").Error
	}); err != nil {
		t.Fatal(err)
	}
	sentinel := errors.New("rollback")
	err := r.WithTransaction(ctx, func(tx *Repository) error {
		if err := tx.db.Exec("INSERT INTO tx_items VALUES (2)").Error; err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v", err)
	}
	if count := itemCount(t, r); count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
}

func TestNestedTransactionUsesSavepoint(t *testing.T) {
	r := transactionRepository(t)
	ctx := context.Background()
	sentinel := errors.New("rollback nested transaction")
	err := r.WithTransaction(ctx, func(tx *Repository) error {
		if err := tx.db.Exec("INSERT INTO tx_items VALUES (1)").Error; err != nil {
			return err
		}
		err := tx.WithTransaction(ctx, func(nested *Repository) error {
			if err := nested.db.Exec("INSERT INTO tx_items VALUES (2)").Error; err != nil {
				return err
			}
			return sentinel
		})
		if !errors.Is(err, sentinel) {
			return errors.New("nested rollback error was lost")
		}
		return tx.db.Exec("INSERT INTO tx_items VALUES (3)").Error
	})
	if err != nil {
		t.Fatal(err)
	}
	if count := itemCount(t, r); count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}
}

func TestCanceledTransactionDoesNotRunCallback(t *testing.T) {
	r := transactionRepository(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	err := r.WithTransaction(ctx, func(*Repository) error { called = true; return nil })
	if called || !errors.Is(err, context.Canceled) {
		t.Fatalf("called=%v error=%v", called, err)
	}
}

func TestTransactionPanicRollsBack(t *testing.T) {
	r := transactionRepository(t)
	func() {
		defer func() {
			if recovered := recover(); recovered != "rollback panic" {
				t.Errorf("panic = %v", recovered)
			}
		}()
		if err := r.WithTransaction(context.Background(), func(tx *Repository) error {
			if err := tx.db.Exec("INSERT INTO tx_items VALUES (1)").Error; err != nil {
				return err
			}
			panic("rollback panic")
		}); err != nil {
			t.Fatal(err)
		}
	}()
	if count := itemCount(t, r); count != 0 {
		t.Fatalf("count = %d, want 0", count)
	}
}
