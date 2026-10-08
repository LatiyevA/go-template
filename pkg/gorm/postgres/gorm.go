package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type config struct {
	translateError  bool
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	connectTimeout  time.Duration
	logMode         logger.LogLevel
	log             *slog.Logger
	nowFunc         func() time.Time
}

// New initializes a bounded pool and verifies connectivity within the timeout.
func New(ctx context.Context, connString string, opts ...Option) (*gorm.DB, error) {
	cfg := &config{
		maxIdleConns:    2,
		maxOpenConns:    10,
		connMaxLifetime: 30 * time.Minute,
		connMaxIdleTime: 5 * time.Minute,
		connectTimeout:  5 * time.Second,
		logMode:         logger.Warn,
		log:             slog.Default(),
		nowFunc:         time.Now,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.maxOpenConns <= 0 || cfg.maxIdleConns < 0 || cfg.maxIdleConns > cfg.maxOpenConns {
		return nil, errors.New("invalid database pool connection limits")
	}
	if cfg.connectTimeout <= 0 || cfg.connMaxLifetime <= 0 || cfg.connMaxIdleTime <= 0 {
		return nil, errors.New("database connection timeouts must be positive")
	}
	db, err := gorm.Open(postgres.Open(connString), cfg.toGormConfig())
	if err != nil {
		return nil, fmt.Errorf("init database session: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("retrieve database pool: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.maxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.connMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.connectTimeout)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, errors.Join(fmt.Errorf("ping database: %w", err), sqlDB.Close())
	}
	return db, nil
}
