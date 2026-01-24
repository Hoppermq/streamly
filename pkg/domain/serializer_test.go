package domain

import (
	"encoding/json"
	"testing"
)

func TestSelectClause_UnmarshalJSON_Field(t *testing.T) {
	t.Parallel()
	jsonData := `"event_name"`

	var sel SelectClause
	if err := json.Unmarshal([]byte(jsonData), &sel); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !sel.IsField() {
		t.Errorf("expected field type, got %s", sel.Type)
	}

	if sel.Field == nil || *sel.Field != "event_name" {
		t.Errorf("expected field='event_name', got %v", sel.Field)
	}
}

func TestSelectClause_UnmarshalJSON_Function(t *testing.T) {
	t.Parallel()
	jsonData := `{"function": "count", "args": ["*"], "alias": "total"}`

	var sel SelectClause
	if err := json.Unmarshal([]byte(jsonData), &sel); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !sel.IsFunction() {
		t.Errorf("expected function type, got %s", sel.Type)
	}

	if sel.Function == nil {
		t.Fatal("function is nil")
	}

	if sel.Function.Function != "count" {
		t.Errorf("expected function='count', got %s", sel.Function.Function)
	}

	if sel.Function.Alias != "total" {
		t.Errorf("expected alias='total', got %s", sel.Function.Alias)
	}
}

func TestQueryAstRequest_UnmarshalJSON_Mixed(t *testing.T) {
	t.Parallel()
	jsonData := `{
		"select": [
			"event_name",
			{"function": "count", "args": ["*"], "alias": "total"}
		],
		"from": "events",
		"time_range": {
			"start": "now-1h",
			"end": "now"
		}
	}`

	var req QueryAstRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(req.Select) != 2 {
		t.Fatalf("expected 2 select clauses, got %d", len(req.Select))
	}

	if !req.Select[0].IsField() {
		t.Errorf("first select should be field")
	}

	if !req.Select[1].IsFunction() {
		t.Errorf("second select should be function")
	}

	if req.From != DataSourceEvents {
		t.Errorf("expected from=events, got %s", req.From)
	}
}

func TestGroupByClause_UnmarshalJSON_Field(t *testing.T) {
	t.Parallel()
	jsonData := `"event_name"`

	var gb GroupByClause
	if err := json.Unmarshal([]byte(jsonData), &gb); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !gb.IsField() {
		t.Errorf("expected field type, got %s", gb.Type)
	}

	if gb.Field == nil || *gb.Field != "event_name" {
		t.Errorf("expected field='event_name', got %v", gb.Field)
	}
}

func TestGroupByClause_UnmarshalJSON_TimeWindow(t *testing.T) {
	t.Parallel()
	jsonData := `{"time_window": "5m", "field": "timestamp"}`

	var gb GroupByClause
	if err := json.Unmarshal([]byte(jsonData), &gb); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !gb.IsTimeWindow() {
		t.Errorf("expected timeWindow type, got %s", gb.Type)
	}

	if gb.TimeWindow == nil {
		t.Fatal("timeWindow is nil")
	}

	if gb.TimeWindow.Window != "5m" {
		t.Errorf("expected window='5m', got %s", gb.TimeWindow.Window)
	}
}
