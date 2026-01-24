package domain

import (
	"encoding/json"

	"github.com/hoppermq/streamly/pkg/domain/errors"
)

func (s *SelectClause) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		s.Type = FieldType
		s.Field = &v
		return nil

	case map[string]any:
		s.Type = FunctionType
		var fn AggFct
		if err := json.Unmarshal(data, &fn); err != nil {
			return errors.SerializerInvalidSelectFunction(err)
		}
		s.Function = &fn
		return nil

	default:
		return errors.ErrSerializerInvalidSelect
	}
}

func (s *SelectClause) MarshalJSON() ([]byte, error) {
	if s.Type == FieldType && s.Field != nil {
		return json.Marshal(*s.Field)
	}
	if s.Type == FunctionType && s.Function != nil {
		return json.Marshal(s.Function)
	}
	return nil, errors.ErrSerializerInvalidSelectClause
}

func (g *GroupByClause) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		g.Type = FieldType
		g.Field = &v
		return nil

	case map[string]any:
		g.Type = TimeWindowType
		var tw TimeWindow
		if err := json.Unmarshal(data, &tw); err != nil {
			return errors.SerializerInvalidTimeWindow(err)
		}
		g.TimeWindow = &tw
		return nil

	default:
		return errors.ErrSerializerInvalidGroupBy
	}
}

func (g *GroupByClause) MarshalJSON() ([]byte, error) {
	if g.Type == FieldType && g.Field != nil {
		return json.Marshal(*g.Field)
	}
	if g.Type == TimeWindowType && g.TimeWindow != nil {
		return json.Marshal(g.TimeWindow)
	}
	return nil, errors.ErrSerializerInvalidGroupByClause
}
