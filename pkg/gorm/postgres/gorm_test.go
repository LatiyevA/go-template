package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestInvalidPoolLimitFailsBeforeConnecting(t *testing.T) {
	if _, err := New(context.Background(), "", MaxOpenConns(0)); err == nil {
		t.Fatal("expected invalid connection limit error")
	}
}

func TestConnectionCheckHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := New(ctx, "postgres://postgres:postgres@localhost:1/app_db?sslmode=disable", SilentLogger())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v", err)
	}
}

func TestDatabasePoolLimits(t *testing.T) {
	url := os.Getenv("TEST_PG_URL")
	if url == "" {
		t.Skip("set TEST_PG_URL to run PostgreSQL integration tests")
	}
	db, err := New(context.Background(), url, MaxOpenConns(3), SilentLogger())
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			t.Error(err)
		}
	}()
	if n := sqlDB.Stats().MaxOpenConnections; n != 3 {
		t.Fatalf("max connections = %d", n)
	}
}
