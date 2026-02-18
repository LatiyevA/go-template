package usecase

import (
	"log/slog"

	"github.com/LatiyevA/go-template/internal/repo"
)

type Usecase struct {
	repo *repo.Repository
	log  *slog.Logger
}

func New(repo *repo.Repository, log *slog.Logger) *Usecase {
	return &Usecase{
		repo: repo,
		log:  log,
	}
}
