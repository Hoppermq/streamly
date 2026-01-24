package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zixyos/glog"
	serviceloader "github.com/zixyos/goloader/service"

	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/core/ingestor"
	"github.com/hoppermq/streamly/internal/core/migration"
	"github.com/hoppermq/streamly/internal/http"
	"github.com/hoppermq/streamly/internal/http/routes"
	"github.com/hoppermq/streamly/internal/storage/clickhouse"
	"github.com/hoppermq/streamly/pkg/domain"
)

//nolint:funlen // ignoring main len.
func main() {
	logger, err := glog.NewDefault()
	if err != nil {
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		).Error("failed to create logger", "error", err)
		os.Exit(domain.ExitStatus)
	}

	ctx := context.Background()

	ingestionConfig, err := config.LoadIngestionConfig()
	logger.Info("configuration loaded", "config", ingestionConfig)
	if err != nil {
		logger.Warn("failed to load ingestion config", "error", err)
	}

	clickhouseDriver := clickhouse.OpenConn(
		clickhouse.WithIngestionConfig(ingestionConfig),
	)

	migrationDriver := clickhouse.OpenConn(
		clickhouse.WithIngestionConfig(ingestionConfig),
	)

	// Extract *sql.DB from Driver for migrations
	var sqlDB *sql.DB
	if chDriver, ok := migrationDriver.(*clickhouse.Driver); ok {
		sqlDB = chDriver.DB()
	}

	migrationService := migration.NewService(
		migration.WithDB(sqlDB),
		migration.WithLogger(logger),
		migration.WithMigrationPath("./clickhouse/sql"),
	)

	if err = migrationService.RunMigrations(ctx); err != nil {
		logger.Error("failed to run migrations", "error", err)
		panic(err)
	}

	eventRepository := ingestor.NewEventRepository(
		ingestor.WithDriver(clickhouseDriver),
	)

	eventUseCase := ingestor.NewEventIngestionUseCase(
		ingestor.UseCaseWithLogger(logger),
		ingestor.WithEventRepository(eventRepository),
	)

	engine := gin.New()

	httpServer := http.NewHTTPServer(
		http.WithEngine(engine),
		http.WithHTTPServer(ingestionConfig),
		http.WithLogger(logger),
		http.WithRoutes(routes.CreateRouteRegistrar(
			routes.CreateIngestionRegistrar(logger, eventUseCase),
		)),
	)

	clickhouseClient := clickhouse.NewClient(
		clickhouse.WithDriver(clickhouseDriver),
		clickhouse.WithLogger(logger),
	)

	ingestionService, err := ingestor.NewIngestor(
		ingestor.WithLogger(logger),
		ingestor.WithConfig(ingestionConfig),
		ingestor.WithHandlers(clickhouseClient, httpServer),
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
