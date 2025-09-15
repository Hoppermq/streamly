-- Rollback migration: Drop all tables in reverse order

DROP TABLE IF EXISTS top_sources_mv;
DROP TABLE IF EXISTS events_hourly_stats_mv;
DROP TABLE IF EXISTS events_minute_stats_mv;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS metrics;
DROP TABLE IF EXISTS errors;
DROP TABLE IF EXISTS streamly_analytics.logs;
DROP DATABASE IF EXISTS streamly_analytics;