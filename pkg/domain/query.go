package domain

import "context"

type Datasource string

const (
	DataSourceEvents              Datasource = "events"
	DataSourceLogs                Datasource = "logs"
	DataSourceErrors              Datasource = "errors"
	DataSourceMetrics             Datasource = "metrics"
	DataSourceEventsMinuteStatsMV Datasource = "events_minute_stats_mv"
	DataSourceEventsHourlyStatsMV Datasource = "events_hourly_stats_mv"
	DataSourceTopSourceMV         Datasource = "top_sources_mv"
)

type AggFct struct {
	Function string   `json:"function"`
	Args     []string `json:"args"`
	Alias    string   `json:"alias"`
}

type TimeWindow struct {
	Window string `json:"time_window"`
	Field  string `json:"field,omitempty"`
}

type SelectClause struct {
	Type     string
	Field    *string
	Function *AggFct
}

type WhereClause struct {
	Field string `json:"field" binding:"required"`
	Op    string `json:"op" binding:"required"`
	Value any    `json:"value" binding:"required"`
}

type GroupByClause struct {
	Type       string
	Field      *string
	TimeWindow *TimeWindow
}

type OrderByClause struct{}

type TimeRange struct {
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

type QueryAstRequest struct {
	Select    []string      `json:"select" binding:"required"`
	From      Datasource    `json:"from" binding:"required"`
	Where     WhereClause   `json:"where,omitempty"`
	GroupBy   GroupByClause `json:"group_by,omitempty"`
	OrderBy   OrderByClause `json:"order_by,omitempty"`
	Limit     uint          `json:"limit" binding:"required"`
	Offset    uint          `json:"offset" binding:"required"`
	TimeRange TimeRange     `json:"time_range" binding:"required"`
}

type QueryRepository interface{}

type QueryUseCase interface {
	SyncQuery(ctx context.Context, req *QueryAstRequest) (any, error)
	AsyncQuery(ctx context.Context, req *QueryAstRequest) (any, error)
}
