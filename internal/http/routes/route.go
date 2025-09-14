package routes

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/internal/http/handlers"
	"github.com/hoppermq/streamly/pkg/domain"
)

// RegisterBaseRoutes will register the routes and group routes to the engine.
func RegisterBaseRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	router.GET("/health", handlers.HealthHandler(func(ctx context.Context) (bool, error) {
		return true, nil
	}))
}

// RegisterIngestionRoutes will register the event ingestion routes.
func RegisterIngestionRoutes(
	router *gin.Engine,
	logger *slog.Logger,

	ingestionUseCase domain.EventIngestionUseCase,
) {
	ingestionHandler := handlers.NewIngestionHandler(
		handlers.WithLogger(logger),
		handlers.WithUSeCase(ingestionUseCase),
	)

	eventsGroup := router.Group("/events")
	{
		eventsGroup.POST("/ingest", ingestionHandler.IngestEvents)
	}
}
