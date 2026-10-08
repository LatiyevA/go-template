package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormLogger struct {
	log   *slog.Logger
	level logger.LogLevel
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &gormLogger{log: l.log, level: level}
}

func (l *gormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Info {
		l.log.InfoContext(ctx, fmt.Sprintf(msg, args...))
	}
}

func (l *gormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Warn {
		l.log.WarnContext(ctx, fmt.Sprintf(msg, args...))
	}
}

func (l *gormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Error {
		l.log.ErrorContext(ctx, fmt.Sprintf(msg, args...))
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.level >= logger.Error:
		_, rows := fc()
		l.log.ErrorContext(ctx, "database query failed", "error", err, "duration", elapsed, "rows", rows)
	case elapsed > 200*time.Millisecond && l.level >= logger.Warn:
		_, rows := fc()
		l.log.WarnContext(ctx, "slow database query", "duration", elapsed, "rows", rows)
	case l.level >= logger.Info:
		_, rows := fc()
		l.log.InfoContext(ctx, "database query", "duration", elapsed, "rows", rows)
	}
}
