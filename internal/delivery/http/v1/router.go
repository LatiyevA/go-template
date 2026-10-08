package v1

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

func addRoutes(h *Handler) http.Handler {
	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		h.log.Error("configure trusted proxies", "error", err)
	}
	router.Use(h.accessLog, h.recoverPanic)
	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
		defer cancel()
		if h.ready == nil || h.ready(ctx) != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	})
	return router
}

func (h *Handler) accessLog(c *gin.Context) {
	start := time.Now()
	c.Next()
	h.log.InfoContext(c.Request.Context(), "HTTP request",
		"method", c.Request.Method, "path", c.Request.URL.Path,
		"status", c.Writer.Status(), "duration", time.Since(start))
}

func (h *Handler) recoverPanic(c *gin.Context) {
	defer func() {
		if recovered := recover(); recovered != nil {
			h.log.ErrorContext(c.Request.Context(), "HTTP handler panic",
				"panic", recovered, "stack", string(debug.Stack()))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	c.Next()
}
