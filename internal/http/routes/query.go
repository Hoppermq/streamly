package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterQueryRoutes will register the query routes.
func RegisterQueryRoutes(
	router *gin.Engine,
	logger *slog.Logger,
) {
	queryGroup := router.Group("/v1")
	{
		queryGroup.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		queryGroup.POST("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

	}
}

// CreateQueryRegistrar returns a RouteRegistrar for query routes.
func CreateQueryRegistrar(logger *slog.Logger) RouteRegistrar {
	return func(engine *gin.Engine) {
		RegisterQueryRoutes(engine, logger)
	}
}
