package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addRoutes(h *Handler) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Implement other routes

	return router
}
