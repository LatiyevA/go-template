package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/LatiyevA/go-template/internal/usecase"
)

type Handler struct {
	uc    *usecase.Usecase
	log   *slog.Logger
	ready func(context.Context) error
	http.Handler
}

func NewHandler(uc *usecase.Usecase, log *slog.Logger, ready func(context.Context) error) *Handler {
	h := &Handler{uc: uc, log: log, ready: ready}
	h.Handler = addRoutes(h)
	return h
}
