package postgres

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/scripts/sql/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type Client struct {
	dbClient *bun.DB
	logger   *slog.Logger
}

type Option func(*Client) error

func WithDB(db *bun.DB) Option {
	return func(c *Client) error {
		c.dbClient = db
		return nil
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil
		}
	}

	return c
}

func (c *Client) Bootstrap(ctx context.Context) error {
	c.logger.InfoContext(ctx, "starting pg boostrap")

	m := migrate.NewMigrations()
	if err := m.Discover(migrations.SqlMigrations); err != nil {
		c.logger.InfoContext(ctx, "failed to discover migrations", "error", err)
		return err
	}

	migrator := migrate.NewMigrator(c.dbClient, m)
	if err := migrator.Init(ctx); err != nil {
		c.logger.WarnContext(ctx, "migrator init failed", "error", err)
		return err
	}

	grp, err := migrator.Migrate(ctx)
	if err != nil {
		c.logger.WarnContext(ctx, "migrator migrate failed", "error", err)
		return err
	}

	if grp.ID == 0 {
		c.logger.Info("no new migrations to run.")
		return nil
	}
	c.logger.InfoContext(ctx, "migrator migration created", "id", grp.ID)

	return nil
}
