// Package clickhouse represent the clickhouse adapter.
package clickhouse

import (
	"context"

	"github.com/hoppermq/streamly/pkg/domain"
)

type Client struct {
	conn domain.Connection
}

type Option func(*Client) error

func WithConn(conn domain.Connection) Option {
	return func(c *Client) error {
		c.conn = conn
		return nil
	}
}

func (c *Client) Run(ctx context.Context) error {
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
