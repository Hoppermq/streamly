package ast

import (
	"encoding/json"
	"log/slog"

	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Validator struct {
	logger  *slog.Logger
	schemas map[string]*jsonschema.Schema
}

type ValidatorOption func(*Validator)

func ValidatorWithLogger(logger *slog.Logger) ValidatorOption {
	return func(v *Validator) {
		v.logger = logger
	}
}

// Execute will execute the ast validation.
func (v *Validator) Execute(data *domain.QueryAstRequest) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var doc any
	if err := json.Unmarshal(jsonBytes, &doc); err != nil {
		return err
	}

	sch := v.schemas["query-ast.schema.json"]
	if err := sch.Validate(doc); err != nil {
		v.logger.Warn("failed to validate schema", "error", err)
		return err
	}

	return nil
}

// RegisterSchema register the schema to the validator registry.
func (v *Validator) RegisterSchema(key string, schema *jsonschema.Schema) {
	v.logger.Info("registering schema", "key", key)
	v.schemas[key] = schema
}

func NewValidator(opts ...ValidatorOption) *Validator {
	v := &Validator{
		schemas: make(map[string]*jsonschema.Schema, 5),
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}
