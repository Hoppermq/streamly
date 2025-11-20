package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/zixyos/glog"
	serviceloader "github.com/zixyos/goloader/service"
)

func main() {
	logger, err := glog.NewDefault() // should prob handle itself logging the error ?
	if err != nil {
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		).Error("failed to initialize logger", "error", err)
		os.Exit(84)
	}

	ctx := context.Background()
	logger.InfoContext(ctx, "starting auth service")
	app := serviceloader.New(serviceloader.WithLogger(logger))

	app.Run(ctx)
}
