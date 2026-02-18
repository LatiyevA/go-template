package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/LatiyevA/go-template/internal/config"
	v1 "github.com/LatiyevA/go-template/internal/delivery/http/v1"
	"github.com/LatiyevA/go-template/internal/repo"
	"github.com/LatiyevA/go-template/internal/usecase"
	"github.com/LatiyevA/go-template/pkg/gorm/postgres"
	"github.com/LatiyevA/go-template/pkg/httpserver"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := postgres.New(cfg.PGURL)
	if err != nil {
		logger.Error(fmt.Sprintf("error connecting to database: %v", err))
		return
	}

	rp := repo.New(db, logger)
	uc := usecase.New(rp, logger)
	handler := v1.NewHandler(uc, logger)
	srv := httpserver.New(handler)
	<-srv.Notify()
}
