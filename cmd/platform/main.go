package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/core/platform"
	"github.com/hoppermq/streamly/internal/core/platform/organization"
	"github.com/hoppermq/streamly/internal/http"
	"github.com/hoppermq/streamly/internal/http/routes"
	"github.com/hoppermq/streamly/internal/storage/postgres"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/zixyos/glog"
	serviceloader "github.com/zixyos/goloader/service"
)

func main() {
	logger, err := glog.NewDefault()
	if err != nil {
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		)
		os.Exit(84)
	}

	ctx := context.Background()
	platformConf, err := config.LoadPlatformConfig()
	if err != nil {
		logger.Warn("failed to load platform config", "error", err)
		os.Exit(84)
	}

	logger.InfoContext(ctx, "starting platform service")

	sqldb := sql.OpenDB(
		pgdriver.NewConnector(pgdriver.WithDSN(platformConf.DatabaseDSN())),
	)

	db := bun.NewDB(sqldb, pgdialect.New())

	d := postgres.NewClient(
		postgres.WithLogger(logger),
		postgres.WithDB(db),
	)

	if err = d.Bootstrap(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to bootstrap database", "error", err)
		os.Exit(84)
	}

	orgRepos, err := organization.NewRepository(
		organization.RepositoryWithLogger(logger),
		organization.RepositoryWithDB(db),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create organization repository", "error", err)
		os.Exit(84)
	}

	generator := uuid.New
	uuidParser := uuid.Parse

	organizationUC, err := organization.NewUseCase(
		organization.UseCaseWithLogger(logger),
		organization.UseCaseWithRepository(orgRepos),
		organization.UseCaseWithGenerator(generator),
		organization.UseCaseWithUUIDParser(uuidParser),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create organization usecase", "error", err)
		os.Exit(84)
	}

	engine := gin.New()
	httpServer := http.NewHTTPServer(
		http.WithEngine(engine),
		http.WithPlatformHTTPServer(platformConf),
		http.WithLogger(logger),
		http.WithRoutes(
			routes.CreateRouteRegistrar(
				routes.CreatePlatformRegistrar(logger, organizationUC),
			),
		),
	)

	platformService := platform.NewStreamlyService(
		platform.WithLogger(logger),
		platform.WithHandler(httpServer),
	)

	app := serviceloader.New(
		serviceloader.WithLogger(logger),
		serviceloader.WithService(platformService),
	)

	app.Run(ctx)
}
