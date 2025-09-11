package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/core/ingestor"
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

	logger.Info("env", os.Getenv("APP_ENV"))

	ingestionConfig, err := config.LoadIngestionConfig()
	logger.Info("conf", ingestionConfig)
	if err != nil {
		logger.Warn("failed to load ingestion config", "error", err)
	}

	ingestionService, err := ingestor.NewIngestor(
		ingestor.WithLogger(logger),
		ingestor.WithConfig(ingestionConfig),
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
