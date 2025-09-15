package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/core/ingestor"
	"github.com/hoppermq/streamly/internal/http"
	"github.com/hoppermq/streamly/internal/storage/clickhouse"
	"github.com/zixyos/glog"
	serviceloader "github.com/zixyos/goloader/service"
)

func main() {
	logger, err := glog.NewDefault()
	if err != nil {
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		).Error("failed to create logger", "error", err)
		os.Exit(84)
	}

	ctx := context.Background()

	ingestionConfig, err := config.LoadIngestionConfig()
	logger.Info("conf", ingestionConfig)
	if err != nil {
		logger.Warn("failed to load ingestion config", "error", err)
	}

	clickhouseDriver := clickhouse.OpenConn(
		clickhouse.WithAddr("http://127.0.0.1:9000"),
		clickhouse.WithUser("admin"),
		clickhouse.WithPassword("admin"),
		clickhouse.WithDatabase("database"),
	)

	eventRepository := ingestor.NewEventRepository(
		ingestor.WithDriver(clickhouseDriver),
	)

	mockEventRepo := ingestor.NewMockEventRepository()
	eventUseCase := ingestor.NewEventIngestionUseCase(
		ingestor.UseCaseWithLogger(logger),
		ingestor.WithEventRepository(mockEventRepo),
		ingestor.WithEventRepository(eventRepository),
	)

	engine := gin.New()

	httpServer := http.NewHTTPServer(
		http.WithEngine(engine),
		http.WithHTTPServer(ingestionConfig),
		http.WithLogger(logger),
		http.WithIngestionUseCase(eventUseCase),
	)

	ingestionService, err := ingestor.NewIngestor(
		ingestor.WithLogger(logger),
		ingestor.WithConfig(ingestionConfig),
		ingestor.WithHandlers(httpServer),
	)

	if err != nil {
		logger.Warn("failed to create ingestion service", "error", err)
		panic(err)
	}

	logger.InfoContext(ctx, "welcome to streamly", "service", "ingestion")
	app := serviceloader.New(
		serviceloader.WithLogger(logger),
		serviceloader.WithService(ingestionService),
	)

	app.Run(ctx)
}
