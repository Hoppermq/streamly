// Package clickhouse represent the clickhouse adapter.
package clickhouse

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/pkg/domain"
)

type Client struct {
	driver domain.Driver

	logger *slog.Logger
}

type Option func(*Client) error

func WithDriver(driver domain.Driver) Option {
	return func(c *Client) error {
		c.driver = driver
		return nil
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

func (c *Client) Run(ctx context.Context) error {
	c.logger.Info("running clickhouse connection")
	go func() {
		tx, err := c.driver.BeginTx(ctx, nil)
		if err != nil {
			c.logger.Warn("failed to begin clickhouse transaction", "error", err)
			return
		}
		err = tx.Commit()
		if err != nil {
			c.logger.Warn("failed to commit clickhouse transaction", "error", err)
			return
		}

	}()

	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	return nil
}

func (c *Client) Name() string {
	return "clickhouse"
}

func (c *Client) IsHealthy() bool {
	return true
}
func NewClient(opts ...Option) *Client {
	c := &Client{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			panic(err)
		}
	}

	return c
}
