package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hoppermq/streamly/internal/http/handlers"
	"github.com/hoppermq/streamly/pkg/domain"
)

// RegisterQueryRoutes will register the query routes.
func RegisterQueryRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	queryUseCase domain.QueryUseCase,
) {
	qh := handlers.NewQueryHandler(
		handlers.WithQueryLogger(logger),
		handlers.WithQueryUseCase(queryUseCase),
	)

	queryGroup := router.Group("/v1")
	{
		queryGroup.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		queryGroup.POST("/queries", qh.Execute)
	}
}

// CreateQueryRegistrar returns a RouteRegistrar for query routes.
func CreateQueryRegistrar(
	logger *slog.Logger,
	queryUseCase domain.QueryUseCase,
) RouteRegistrar {
	return func(engine *gin.Engine) {
		RegisterQueryRoutes(engine, logger, queryUseCase)
	}
}
