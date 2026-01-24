package testcontainers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	*redis.RedisContainer
	ConnectionString string
}

func StartRedis(ctx context.Context) (*RedisContainer, error) {
	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	host, err := redisContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis host: %w", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get redis port: %w", err)
	}

	connStr := net.JoinHostPort(host, port.Port())

	return &RedisContainer{
		RedisContainer:   redisContainer,
		ConnectionString: connStr,
	}, nil
}

func (c *RedisContainer) Close(ctx context.Context) error {
	return c.Terminate(ctx)
}
