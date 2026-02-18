package repo

import (
	"log/slog"

	"gorm.io/gorm"
)

type (
	Repository struct {
		db  *gorm.DB
		log *slog.Logger
	}
)

func New(db *gorm.DB, log *slog.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}
