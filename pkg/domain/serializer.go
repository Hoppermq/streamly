package domain

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (s *SelectClause) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		s.Type = "field"
		s.Field = &v
		return nil

	case map[string]any:
		s.Type = "function"
		var fn AggFct
		if err := json.Unmarshal(data, &fn); err != nil {
			return fmt.Errorf("invalid function select: %w", err)
		}
		s.Function = &fn
		return nil

	default:
		return errors.New("select must be string or aggregation object")
	}
}

func (s SelectClause) MarshalJSON() ([]byte, error) {
	if s.Type == "field" && s.Field != nil {
		return json.Marshal(*s.Field)
	}
	if s.Type == "function" && s.Function != nil {
		return json.Marshal(s.Function)
	}
	return nil, errors.New("invalid select clause")
}

func (g *GroupByClause) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		g.Type = "field"
		g.Field = &v
		return nil

	case map[string]any:
		g.Type = "timeWindow"
		var tw TimeWindow
		if err := json.Unmarshal(data, &tw); err != nil {
			return fmt.Errorf("invalid time window: %w", err)
		}
		g.TimeWindow = &tw
		return nil

	default:
		return errors.New("groupBy must be string or time window object")
	}
}

func (g GroupByClause) MarshalJSON() ([]byte, error) {
	if g.Type == "field" && g.Field != nil {
		return json.Marshal(*g.Field)
	}
	if g.Type == "timeWindow" && g.TimeWindow != nil {
		return json.Marshal(g.TimeWindow)
	}
	return nil, errors.New("invalid groupBy clause")
}
