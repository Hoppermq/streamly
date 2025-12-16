package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/storage/postgres"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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

	for {
	}
}
