package postgres

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Option func(*config)

func NowFunc(f func() time.Time) Option {
	return func(c *config) { c.nowFunc = f }
}

func TranslateError(enable bool) Option {
	return func(c *config) { c.translateError = enable }
}

func MaxIdleConns(n int) Option {
	return func(c *config) { c.maxIdleConns = n }
}

func MaxOpenConns(n int) Option {
	return func(c *config) { c.maxOpenConns = n }
}

func ConnMaxLifetime(d time.Duration) Option {
	return func(c *config) { c.connMaxLifetime = d }
}

func ConnMaxIdleTime(d time.Duration) Option {
	return func(c *config) { c.connMaxIdleTime = d }
}

func ConnectTimeout(d time.Duration) Option {
	return func(c *config) { c.connectTimeout = d }
}

func Logger(log *slog.Logger) Option {
	return func(c *config) { c.log = log }
}

func SilentLogger() Option {
	return func(c *config) { c.logMode = logger.Silent }
}

func (c *config) toGormConfig() *gorm.Config {
	return &gorm.Config{
		NowFunc:        c.nowFunc,
		TranslateError: c.translateError,
		// Ping explicitly with a context instead of GORM's unbounded automatic ping.
		DisableAutomaticPing: true,
		Logger:               &gormLogger{log: c.log, level: c.logMode},
	}
}
