package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/core/auth"
	"github.com/hoppermq/streamly/internal/http"
	"github.com/hoppermq/streamly/internal/http/routes"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/zixyos/glog"
	serviceloader "github.com/zixyos/goloader/service"
)

func main() {
	logger, err := glog.NewDefault() // should prob handle itself logging the error ?
	if err != nil {
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		).Error("failed to initialize logger", "error", err)
		os.Exit(domain.ExitStatus)
	}

	ctx := context.Background()
	authConf, err := config.LoadAuthConfig()
	if err != nil {
		logger.Warn("failed to load auth config", "error", err)
		// TODO: default config should be run
	}

	logger.InfoContext(ctx, "starting auth service")

	engine := gin.New()

	httpServer := http.NewHTTPServer(
		http.WithEngine(engine),
		http.WithAuthHTTPServer(authConf),
		http.WithLogger(logger),
		http.WithRoutes(routes.CreateRouteRegistrar(
			routes.CreateAuthRegistrar(logger),
		)),
	)

	authService := auth.NewAuthService(
		auth.WithLogger(logger),
		auth.WithHandler(httpServer),
	)

	app := serviceloader.New(
		serviceloader.WithLogger(logger),
		serviceloader.WithService(authService),
	)

	app.Run(ctx)
}
