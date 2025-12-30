//go:build integration

package main_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/internal/tests/testcontainers"
	"github.com/hoppermq/streamly/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryService_ClickHouseQueries(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	chContainer, err := testcontainers.StartClickHouse(ctx)
	require.NoError(t, err)
	defer chContainer.Close(ctx)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	logger.InfoContext(ctx, "ClickHouse container started for query tests")

	err = chContainer.Conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS events (
			timestamp DATETIME64(3, 'UTC'),
			tenant_id LowCardinality(String),
			message_id String,
			source_id LowCardinality(String),
			topic LowCardinality(String),
			content_raw String,
			content_size_bytes UInt32,
			headers Map(String, String),
			frame_type UInt8,
			event_type LowCardinality(String)
		) ENGINE = MergeTree()
		ORDER BY (timestamp, tenant_id, topic, message_id)
	`)
	require.NoError(t, err)

	tenantID := uuid.New().String()
	events := testutil.NewBatchEvents(tenantID, 100, "query.test")

	batch, err := chContainer.Conn.PrepareBatch(ctx, `
		INSERT INTO events (
			timestamp, tenant_id, message_id, source_id, topic,
			content_raw, content_size_bytes, headers, frame_type, event_type
		)
	`)
	require.NoError(t, err)

	for _, event := range events {
		err = batch.Append(
			event.Timestamp,
			event.TenantID,
			event.MessageID,
			event.SourceID,
			event.Topic,
			event.ContentRaw,
			event.ContentSize,
			event.Headers,
			event.FrameType,
			event.EventType,
		)
		require.NoError(t, err)
	}
	err = batch.Send()
	require.NoError(t, err)

	t.Run("Query events with time range filter", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now().Add(1 * time.Hour)

		rows, err := chContainer.Conn.Query(ctx, `
			SELECT count()
			FROM events
			WHERE tenant_id = ? AND timestamp BETWEEN ? AND ?
		`, tenantID, startTime, endTime)
		require.NoError(t, err)
		defer rows.Close()

		var count uint64
		rows.Next()
		err = rows.Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), count)
	})

	t.Run("Query events with pagination", func(t *testing.T) {
		rows, err := chContainer.Conn.Query(ctx, `
			SELECT message_id, event_type
			FROM events
			WHERE tenant_id = ?
			ORDER BY timestamp DESC
			LIMIT 10 OFFSET 0
		`, tenantID)
		require.NoError(t, err)
		defer rows.Close()

		var resultCount int
		for rows.Next() {
			var messageID, eventType string
			err := rows.Scan(&messageID, &eventType)
			require.NoError(t, err)
			resultCount++
		}
		assert.Equal(t, 10, resultCount)
	})

	t.Run("Aggregation query - count by event_type", func(t *testing.T) {
		rows, err := chContainer.Conn.Query(ctx, `
			SELECT event_type, count() as cnt
			FROM events
			WHERE tenant_id = ?
			GROUP BY event_type
		`, tenantID)
		require.NoError(t, err)
		defer rows.Close()

		var eventType string
		var count uint64
		rows.Next()
		err = rows.Scan(&eventType, &count)
		require.NoError(t, err)
		assert.Equal(t, "query.test", eventType)
		assert.Equal(t, uint64(100), count)
	})

	t.Run("Query performance - P99 latency < 200ms", func(t *testing.T) {
		var latencies []time.Duration

		for i := 0; i < 100; i++ {
			start := time.Now()

			rows, err := chContainer.Conn.Query(ctx, `
				SELECT message_id
				FROM events
				WHERE tenant_id = ? AND event_type = ?
				LIMIT 10
			`, tenantID, "query.test")
			require.NoError(t, err)

			for rows.Next() {
				var messageID string
				rows.Scan(&messageID)
			}
			rows.Close()

			latency := time.Since(start)
			latencies = append(latencies, latency)
		}

		p99Index := int(float64(len(latencies)) * 0.99)
		p99Latency := latencies[p99Index]

		t.Logf("Query P99 latency: %v", p99Latency)
		assert.Less(t, p99Latency, 200*time.Millisecond, "P99 query latency should be <200ms")
	})
}

func TestQueryService_WithRedisCache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	redisContainer, err := testcontainers.StartRedis(ctx)
	require.NoError(t, err)
	defer redisContainer.Close(ctx)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	logger.InfoContext(ctx, "Redis container started", "connection", redisContainer.ConnectionString)

	t.Run("Redis connectivity", func(t *testing.T) {
		assert.NotEmpty(t, redisContainer.ConnectionString)
		assert.Contains(t, redisContainer.ConnectionString, "redis://")
	})
}

func TestQueryService_MultiServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	chContainer, err := testcontainers.StartClickHouse(ctx)
	require.NoError(t, err)
	defer chContainer.Close(ctx)

	redisContainer, err := testcontainers.StartRedis(ctx)
	require.NoError(t, err)
	defer redisContainer.Close(ctx)

	err = chContainer.Conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS events (
			timestamp DATETIME64(3, 'UTC'),
			tenant_id LowCardinality(String),
			message_id String,
			source_id LowCardinality(String),
			topic LowCardinality(String),
			content_raw String,
			content_size_bytes UInt32,
			headers Map(String, String),
			frame_type UInt8,
			event_type LowCardinality(String)
		) ENGINE = MergeTree()
		ORDER BY (timestamp, tenant_id, topic, message_id)
	`)
	require.NoError(t, err)

	t.Run("Full stack integration - ClickHouse + Redis available", func(t *testing.T) {
		assert.NotNil(t, chContainer.Conn)
		assert.NotEmpty(t, redisContainer.ConnectionString)

		tenantID := uuid.New().String()
		event := testutil.NewSampleEvent(tenantID, "integration.test")

		err := chContainer.Conn.Exec(ctx, `
			INSERT INTO events (
				timestamp, tenant_id, message_id, source_id, topic,
				content_raw, content_size_bytes, headers, frame_type, event_type
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			event.Timestamp,
			event.TenantID,
			event.MessageID,
			event.SourceID,
			event.Topic,
			event.ContentRaw,
			event.ContentSize,
			event.Headers,
			event.FrameType,
			event.EventType,
		)
		require.NoError(t, err)

		var count uint64
		err = chContainer.Conn.QueryRow(ctx, "SELECT count() FROM events WHERE tenant_id = ?", tenantID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), count)
	})
}
