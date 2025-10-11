package ast

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hoppermq/streamly/internal/storage/clickhouse"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Builder struct {
	logger       *slog.Logger
	schemaFS     embed.FS
	jschCompiler *jsonschema.Compiler
	validator    *Validator
	translator   *clickhouse.Translator // should be domain driver.
}

type BuilderOption func(builder *Builder)

func BuilderWithLogger(logger *slog.Logger) BuilderOption {
	return func(builder *Builder) {
		builder.logger = logger
	}
}

func BuilderWithValidator(validator *Validator) BuilderOption {
	return func(builder *Builder) {
		builder.validator = validator
	}
}

func BuilderWithSchemaFS(schemaFS embed.FS) BuilderOption {
	return func(translator *Builder) {
		translator.schemaFS = schemaFS
	}
}

func BuilderWithJsonSchemaCompiler(jsonSchemaCompiler *jsonschema.Compiler) BuilderOption {
	return func(builder *Builder) {
		builder.jschCompiler = jsonSchemaCompiler
	}
}

func BuilderWithTranslator(translator *clickhouse.Translator) BuilderOption {
	return func(builder *Builder) {
		builder.translator = translator
	}
}

func (tr *Builder) Run(ctx context.Context) error {
	schemaContent, err := tr.schemaFS.ReadFile("query-ast.schema.json")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	var schDoc any
	err = json.Unmarshal(schemaContent, &schDoc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	if err := tr.jschCompiler.AddResource("query-ast.schema.json", schDoc); err != nil {
		return fmt.Errorf("failed to add schema resource: %w", err)
	}

	schema, err := tr.jschCompiler.Compile("query-ast.schema.json")
	if err != nil {
		return fmt.Errorf("failed to compile schema: %w", err)
	}

	if schema == nil {
		tr.logger.Warn("schema is nil, skipping...")
		return fmt.Errorf("failed to compile schema")
	}
	tr.validator.RegisterSchema("query-ast.schema.json", schema)

	return nil
}

func (tr *Builder) Shutdown(ctx context.Context) error {
	tr.logger.Info("shutting down ast builder handler")
	return nil
}

func (tr *Builder) Name() string {
	return "ast-builder-handler"
}

func (tr *Builder) IsHealthy() bool {
	return true
}

func (tr *Builder) Execute(data *domain.QueryAstRequest) (domain.Query, []domain.QueryArgs, error) { // should be a domain type here
	if err := tr.validator.Execute(data); err != nil {
		tr.logger.Warn("error while executing the validation")
		return "", nil, err
	}
	query, err := tr.translator.Translate(data)
	if err != nil {
		tr.logger.Warn("error while translating data")
		return "", nil, err
	}
	q, params, err := query.Build()
	if err != nil {
		tr.logger.Warn("error while building query")
		return "", nil, err
	}

	queryArgs := make([]domain.QueryArgs, len(params))
	for i, param := range params {
		queryArgs[i] = param
	}

	tr.logger.Info("query built", "query", q, "params", params)
	return domain.Query(q), queryArgs, nil
}

func NewBuilder(opts ...BuilderOption) *Builder {
	translator := &Builder{}
	for _, opt := range opts {
		opt(translator)
	}

	return translator
}
