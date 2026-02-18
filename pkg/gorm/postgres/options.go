package postgres

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Option func(*config)

func NowFunc(f func() time.Time) Option {
	return func(c *config) {
		c.nowFunc = f
	}
}

func TranslateError(enable bool) Option {
	return func(c *config) {
		c.translateError = enable
	}
}

func MaxIdleConns(n int) Option {
	return func(c *config) {
		c.maxIdleConns = n
	}
}

func MaxOpenConns(n int) Option {
	return func(c *config) {
		c.maxOpenConns = n
	}
}

func SilentLogger() Option {
	return func(c *config) {
		c.logMode = logger.Silent
	}
}

func (c *config) toGormConfig() *gorm.Config {
	return &gorm.Config{
		NowFunc:        c.nowFunc,
		TranslateError: c.translateError,
		Logger:         logger.Default.LogMode(c.logMode),
	}
}
