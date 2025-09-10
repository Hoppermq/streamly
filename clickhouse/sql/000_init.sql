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

CREATE TABLE IF NOT EXISTS events(
  timestamp DATETIME('Europe/Paris') CODEC(Delta, LZ4),
  tenant_id LowCardinality(String),

  message_id String CODEC(ZSTD),
  source_id LowCardinality(String),
  topic LowCardinality(String),

  content_raw String Codec(ZSTD),
  content_json JSON,
  content_size_bytes UInt32,

  headers Map(String, String) CODEC(ZSTD),

  frame_type UInt8,

  event_type LowCardinality(String)
) ENGINE = MergeTree()
PARTITION BY (toYYYYMM(timestamp), tenant_id, topic)
ORDER BY (timestamp, tenant_id, topic, message_id);

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
