package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/LatiyevA/go-template/internal/config"
	v1 "github.com/LatiyevA/go-template/internal/delivery/http/v1"
	"github.com/LatiyevA/go-template/internal/repo"
	"github.com/LatiyevA/go-template/internal/usecase"
	"github.com/LatiyevA/go-template/pkg/gorm/postgres"
	"github.com/LatiyevA/go-template/pkg/httpserver"
	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err := run(logger); err != nil {
		logger.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) (err error) {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	logger = logger.With("service", cfg.ServiceName)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := postgres.New(ctx, cfg.PGURL,
		postgres.MaxOpenConns(cfg.MaxOpenConns),
		postgres.MaxIdleConns(cfg.MaxIdleConns),
		postgres.ConnMaxLifetime(cfg.ConnMaxLifetime),
		postgres.ConnMaxIdleTime(cfg.ConnMaxIdleTime),
		postgres.ConnectTimeout(cfg.ConnectTimeout),
		postgres.Logger(logger),
	)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("retrieve database pool: %w", err)
	}
	defer func() { err = errors.Join(err, sqlDB.Close()) }()

	rp := repo.New(db, logger)
	uc := usecase.New(rp, logger)
	handler := v1.NewHandler(uc, logger, sqlDB.PingContext)
	srv := httpserver.New(handler,
		httpserver.Port(strconv.Itoa(cfg.HTTPPort)),
		httpserver.ShutdownTimeout(cfg.ShutdownTimeout),
		httpserver.Logger(logger),
	)
	if err := srv.Start(); err != nil {
		return fmt.Errorf("start HTTP server: %w", err)
	}
	logger.Info("HTTP server started", "port", cfg.HTTPPort)

	select {
	case <-ctx.Done():
		stop()
		logger.Info("shutting down HTTP server")
		if err := srv.Shutdown(); err != nil {
			return errors.Join(fmt.Errorf("shutdown HTTP server: %w", err), srv.Close())
		}
		return nil
	case err := <-srv.Notify():
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return errors.Join(fmt.Errorf("serve HTTP: %w", err), srv.Close())
	}
}
