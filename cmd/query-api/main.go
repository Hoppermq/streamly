package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/core/query"
	"github.com/hoppermq/streamly/internal/http"
	"github.com/hoppermq/streamly/internal/http/routes"
	"github.com/zixyos/glog"
	serviceloader "github.com/zixyos/goloader/service"
)

func main() {
	logger, err := glog.NewDefault()
	if err != nil {
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		).Error("failed to initialize logger", "error", err)
		os.Exit(84)
	}

	ctx := context.Background()

	queryConfig, err := config.LoadQueryConfig()
	if err != nil {
		logger.Warn("failed to load query config", "error", err)
	}

	engine := gin.New()

	queryUseCase := query.NewQueryUseCase(
		query.UseCaseWithLogger(logger),
	)

	httpServer := http.NewHTTPServer(
		http.WithEngine(engine),
		http.WithQueryHTTPServer(queryConfig),
		http.WithLogger(logger),
		http.WithRoutes(
			routes.CreateRouteRegistrar(
				routes.CreateQueryRegistrar(logger, queryUseCase),
			),
		),
	)

	queryService := query.NewQueryService(
		query.WithLogger(logger),
		query.WithHandlers(httpServer),
	)

	app := serviceloader.New(
		serviceloader.WithLogger(logger),
		serviceloader.WithService(queryService),
	)

	app.Run(ctx)
}
