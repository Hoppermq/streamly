package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/hoppermq/streamly/internal/storage/postgres"
	"github.com/zixyos/glog"
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
	logger.InfoContext(ctx, "starting platform service")

	db, err := sql.Open("pg", "")
	if err != nil {
		logger.ErrorContext(ctx, "failed to open database", err)
		os.Exit(84)
	}

	d := postgres.NewClient(
		postgres.WithLogger(logger),
		postgres.WithDB(db),
	)

	if err = d.Bootstrap(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to bootstrap database", err)
		os.Exit(84)
	}
}
