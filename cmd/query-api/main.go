package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/hoppermq/streamly/internal/core/query"
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

	queryService := query.NewQueryService(
		query.WithLogger(logger),
	)

	app := serviceloader.New(
		serviceloader.WithLogger(logger),
		serviceloader.WithService(queryService))

	app.Run(ctx)
}
