package testcontainers

import (
	"context"
	"fmt"
	"net"
	"time"

	clickhouseDriver "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/clickhouse"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ClickHouseContainer struct {
	*clickhouse.ClickHouseContainer
	ConnectionString string
	Conn             clickhouseDriver.Conn
}

func StartClickHouse(ctx context.Context) (*ClickHouseContainer, error) {
	chContainer, err := clickhouse.Run(ctx,
		"clickhouse/clickhouse-server:24.1-alpine",
		testcontainers.WithEnv(map[string]string{
			"CLICKHOUSE_DB":       "default",
			"CLICKHOUSE_USER":     "default",
			"CLICKHOUSE_PASSWORD": "",
		}),
		testcontainers.WithWaitStrategy(
			wait.ForHTTP("/ping").
				WithPort("8123/tcp").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start clickhouse container: %w", err)
	}

	host, err := chContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get clickhouse host: %w", err)
	}

	port, err := chContainer.MappedPort(ctx, "9000/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get clickhouse port: %w", err)
	}

	connStr := net.JoinHostPort("clickhouse://"+host, port.Port())

	conn, err := clickhouseDriver.Open(&clickhouseDriver.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port.Port())},
		Auth: clickhouseDriver.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clickhouse: %w", err)
	}

	return &ClickHouseContainer{
		ClickHouseContainer: chContainer,
		ConnectionString:    connStr,
		Conn:                conn,
	}, nil
}

func (c *ClickHouseContainer) Close(ctx context.Context) error {
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			return err
		}
	}
	return c.Terminate(ctx)
}
