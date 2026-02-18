package postgres

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	_defaultMaxIdleConns = 2
	_defaultMaxOpenConns = 0
	_defaultLogMode      = logger.Info
)

type config struct {
	translateError bool
	maxIdleConns   int
	maxOpenConns   int
	logMode        logger.LogLevel
	nowFunc        func() time.Time
}

func New(connString string, opts ...Option) (*gorm.DB, error) {
	defaultNowFunc := func() time.Time { return time.Now() }

	cfg := &config{
		maxIdleConns: _defaultMaxIdleConns,
		maxOpenConns: _defaultMaxOpenConns,
		logMode:      _defaultLogMode,
		nowFunc:      defaultNowFunc,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	db, err := gorm.Open(postgres.Open(connString), cfg.toGormConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to init db session: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sqlDB object: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.maxOpenConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
