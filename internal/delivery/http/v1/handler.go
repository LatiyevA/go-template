package v1

import (
	"log/slog"
	"net/http"

	"github.com/LatiyevA/go-template/internal/usecase"
)

type Handler struct {
	uc  *usecase.Usecase
	log *slog.Logger

	http.Handler
}

func NewHandler(uc *usecase.Usecase, log *slog.Logger) *Handler {

	h := &Handler{
		uc:  uc,
		log: log,
	}

	h.Handler = addRoutes(h)
	return h
}
