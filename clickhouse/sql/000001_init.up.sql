CREATE DATABASE IF NOT EXISTS streamly_analytics;

CREATE TABLE IF NOT EXISTS streamly_analytics.logs (
  timestamp DATETIME('Europe/Paris') CODEC(Delta, LZ4),
  tenant_id LowCardinality(String),

  metric_name LowCardinality(String),
  service_name LowCardinality(String),
  instance LowCardinality(String),

  value Float64 CODEC(Gorilla, LZ4),

  labels Map(String, String),

  data_size_bytes UInt32
) ENGINE = MergeTree()
PARTITION BY (toYYYYMM(timestamp), tenant_id)
ORDER BY (tenant_id, metric_name, timestamp);

CREATE TABLE IF NOT EXISTS events (
    -- Core event data (matches your domain.Event exactly)
    timestamp DATETIME64(3, 'UTC') CODEC(Delta, LZ4),
    tenant_id LowCardinality(String),
    message_id String CODEC(ZSTD),
    source_id LowCardinality(String),
    topic LowCardinality(String),

    -- Content fields
    content_raw String CODEC(ZSTD),
    content_json JSON,
    content_size_bytes UInt32,

    -- Metadata
    headers Map(String, String) CODEC(ZSTD),
    frame_type UInt8,
    event_type LowCardinality(String),

    -- Analytics optimization (computed on insert)
    ingestion_timestamp DATETIME64(3, 'UTC') DEFAULT now64(),
    date_partition Date MATERIALIZED toDate(timestamp),
    hour_bucket DATETIME MATERIALIZED toStartOfHour(timestamp)

) ENGINE = MergeTree()
PARTITION BY (date_partition, tenant_id, topic)
ORDER BY (timestamp, tenant_id, topic, message_id)
SETTINGS
    index_granularity = 8192,
    merge_with_ttl_timeout = 3600;

CREATE TABLE IF NOT EXISTS errors(
  timestamp DATETIME('Europe/Paris') CODEC(Delta, LZ4HC),
  tenant_id LowCardinality(String)
) ENGINE = MergeTree()
PARTITION BY (toYYYYMM(timestamp), tenant_id)
ORDER BY (timestamp, tenant_id);

CREATE TABLE IF NOT EXISTS metrics(
  timestamp DATETIME('Europe/Paris') CODEC(Delta, LZ4),
  tenant_id LowCardinality(String)
) ENGINE = MergeTree()
PARTITION BY (toYYYYMM(timestamp), tenant_id)
ORDER BY (timestamp, tenant_id);

-- Analytics Materialized Views for Real-time Dashboards

-- 1. Event counts by minute (for real-time charts)
CREATE MATERIALIZED VIEW IF NOT EXISTS events_minute_stats_mv
ENGINE = SummingMergeTree()
PARTITION BY (toYYYYMM(minute_bucket), tenant_id)
ORDER BY (tenant_id, topic, event_type, minute_bucket)
AS SELECT
    tenant_id,
    topic,
    event_type,
    source_id,
    toStartOfMinute(timestamp) as minute_bucket,
    count() as event_count,
    sum(content_size_bytes) as total_bytes,
    uniq(message_id) as unique_messages
FROM events
GROUP BY tenant_id, topic, event_type, source_id, minute_bucket;

-- 2. Hourly aggregations (for historical analysis)
CREATE MATERIALIZED VIEW IF NOT EXISTS events_hourly_stats_mv
ENGINE = SummingMergeTree()
PARTITION BY (toYYYYMM(hour_bucket), tenant_id)
ORDER BY (tenant_id, topic, hour_bucket)
AS SELECT
    tenant_id,
    topic,
    event_type,
    toStartOfHour(timestamp) as hour_bucket,
    count() as event_count,
    sum(content_size_bytes) as total_bytes,
    uniq(source_id) as unique_sources,
    min(timestamp) as first_event,
    max(timestamp) as last_event
FROM events
GROUP BY tenant_id, topic, event_type, hour_bucket;

-- 3. Top sources view (for dashboard "top producers")
CREATE MATERIALIZED VIEW IF NOT EXISTS top_sources_mv
ENGINE = ReplacingMergeTree()
PARTITION BY (toYYYYMM(hour_bucket), tenant_id)
ORDER BY (tenant_id, hour_bucket, source_id)
AS SELECT
    tenant_id,
    source_id,
    toStartOfHour(timestamp) as hour_bucket,
    count() as event_count,
    sum(content_size_bytes) as total_bytes,
    uniqArray(groupArray(topic)) as topics_used
FROM events
GROUP BY tenant_id, source_id, hour_bucket;
