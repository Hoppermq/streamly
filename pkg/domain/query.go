package domain

import "context"

type Datasource string
type Query string
type QueryArgs any

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
	Alias    string   `json:"alias,omitempty"`
}

type TimeWindow struct {
	Window string `json:"time_window"`
	Field  string `json:"field,omitempty"`
}

type Clause string

const (
	FieldType      Clause = "field"
	FunctionType   Clause = "function"
	TimeWindowType Clause = "time_window"
)

type SelectClause struct {
	Type     Clause
	Field    *string
	Function *AggFct
}

func (s *SelectClause) IsField() bool    { return s.Type == "field" }
func (s *SelectClause) IsFunction() bool { return s.Type == "function" }

type WhereClause struct {
	Field string `json:"field"`
	Op    string `json:"op"`
	Value any    `json:"value"`
}

type GroupByClause struct {
	Type       Clause
	Field      *string
	TimeWindow *TimeWindow
}

func (g *GroupByClause) IsField() bool      { return g.Type == FieldType }
func (g *GroupByClause) IsTimeWindow() bool { return g.Type == TimeWindowType }

type OrderByClause struct {
	Field     string `json:"field"`
	Direction string `json:"direction,omitempty"`
}

type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type QueryAstRequest struct {
	Select    []SelectClause  `json:"select"`
	From      Datasource      `json:"from"`
	TimeRange TimeRange       `json:"time_range"`
	Where     []WhereClause   `json:"where,omitempty"`
	GroupBy   []GroupByClause `json:"group_by,omitempty"`
	OrderBy   []OrderByClause `json:"order_by,omitempty"`
	Limit     *int            `json:"limit,omitempty"`
	Offset    *int            `json:"offset,omitempty"`

	TenantID  string `json:"-"`
	RequestID string `json:"-"`
}

type QueryResponse struct {
	RequestID string           `json:"request_id"`
	Data      []map[string]any `json:"data"`
	RowCount  int              `json:"row_count"`
}

type QueryRepository interface {
	ExecuteQuery(ctx context.Context, query Query, args ...QueryArgs) (*QueryResponse, error)
}

type QueryUseCase interface {
	SyncQuery(ctx context.Context, req *QueryAstRequest) (*QueryResponse, error)
}
