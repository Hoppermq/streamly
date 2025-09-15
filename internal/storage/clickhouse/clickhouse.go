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
		_, err := c.driver.Begin()
		if err != nil {
			panic(err)
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
