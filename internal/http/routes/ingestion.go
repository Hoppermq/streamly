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
