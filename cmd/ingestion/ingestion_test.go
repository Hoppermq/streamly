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

func TestIngestionService_ClickHouseIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	chContainer, err := testcontainers.StartClickHouse(ctx)
	require.NoError(t, err)
	defer chContainer.Close(ctx)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	logger.InfoContext(ctx, "ClickHouse container started", "connection", chContainer.ConnectionString)

	t.Run("Create events table schema", func(t *testing.T) {
		err := chContainer.Conn.Exec(ctx, `
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
				event_type LowCardinality(String),
				ingestion_timestamp DATETIME64(3, 'UTC') DEFAULT now64()
			) ENGINE = MergeTree()
			ORDER BY (timestamp, tenant_id, topic, message_id)
		`)
		require.NoError(t, err)
	})

	t.Run("Insert single event into ClickHouse", func(t *testing.T) {
		tenantID := uuid.New().String()
		event := testutil.NewSampleEvent(tenantID, "user.created")

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

	t.Run("Batch insert 100 events", func(t *testing.T) {
		tenantID := uuid.New().String()
		events := testutil.NewBatchEvents(tenantID, 100, "batch.test")

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

		var count uint64
		err = chContainer.Conn.QueryRow(ctx, "SELECT count() FROM events WHERE tenant_id = ?", tenantID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), count)
	})

	t.Run("Query events by tenant_id and event_type", func(t *testing.T) {
		tenantID := uuid.New().String()
		events := testutil.NewBatchEvents(tenantID, 50, "query.test")

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

		rows, err := chContainer.Conn.Query(ctx, `
			SELECT message_id, event_type, timestamp
			FROM events
			WHERE tenant_id = ? AND event_type = ?
			ORDER BY timestamp DESC
			LIMIT 10
		`, tenantID, "query.test")
		require.NoError(t, err)
		defer rows.Close()

		var resultCount int
		for rows.Next() {
			var messageID, eventType string
			var timestamp time.Time
			err := rows.Scan(&messageID, &eventType, &timestamp)
			require.NoError(t, err)
			assert.Equal(t, "query.test", eventType)
			resultCount++
		}
		assert.Equal(t, 10, resultCount)
	})

	t.Run("Verify tenant isolation", func(t *testing.T) {
		tenant1 := uuid.New().String()
		tenant2 := uuid.New().String()

		events1 := testutil.NewBatchEvents(tenant1, 10, "tenant1.event")
		events2 := testutil.NewBatchEvents(tenant2, 15, "tenant2.event")

		batch, err := chContainer.Conn.PrepareBatch(ctx, `
			INSERT INTO events (
				timestamp, tenant_id, message_id, source_id, topic,
				content_raw, content_size_bytes, headers, frame_type, event_type
			)
		`)
		require.NoError(t, err)

		for _, event := range append(events1, events2...) {
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

		var count1 uint64
		err = chContainer.Conn.QueryRow(ctx, "SELECT count() FROM events WHERE tenant_id = ?", tenant1).Scan(&count1)
		require.NoError(t, err)
		assert.Equal(t, uint64(10), count1)

		var count2 uint64
		err = chContainer.Conn.QueryRow(ctx, "SELECT count() FROM events WHERE tenant_id = ?", tenant2).Scan(&count2)
		require.NoError(t, err)
		assert.Equal(t, uint64(15), count2)
	})
}

func TestIngestionService_PerformanceBaseline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	chContainer, err := testcontainers.StartClickHouse(ctx)
	require.NoError(t, err)
	defer chContainer.Close(ctx)

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

	t.Run("Insert 1000 events performance test", func(t *testing.T) {
		tenantID := uuid.New().String()
		events := testutil.NewBatchEvents(tenantID, 1000, "perf.test")

		start := time.Now()

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

		duration := time.Since(start)
		eventsPerSecond := float64(1000) / duration.Seconds()

		t.Logf("Inserted 1000 events in %v (%.0f events/sec)", duration, eventsPerSecond)
		assert.Less(t, duration, 5*time.Second, "Batch insert should complete in <5s")
	})
}
