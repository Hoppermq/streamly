package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/internal/http/handlers"
	"github.com/hoppermq/streamly/pkg/domain"
)

// RegisterIngestionRoutes will register the event ingestion routes.
func RegisterIngestionRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	ingestionUseCase domain.IngestionUseCase,
) {
	ingestionHandler := handlers.NewIngestionHandler(
		handlers.WithLogger(logger),
		handlers.WithUSeCase(ingestionUseCase),
	)

	v1 := router.Group("/v1")
	eventsGroup := v1.Group("/events")
	{
		eventsGroup.POST("/ingest", ingestionHandler.IngestEvents)
	}
}

// CreateIngestionRegistrar returns a RouteRegistrar for ingestion routes.
func CreateIngestionRegistrar(
	logger *slog.Logger,
	ingestionUseCase domain.IngestionUseCase,
) RouteRegistrar {
	return func(engine *gin.Engine) {
		RegisterIngestionRoutes(engine, logger, ingestionUseCase)
	}
}
