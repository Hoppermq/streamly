package ast

import (
	"log/slog"
)

type Builder struct {
	logger *slog.Logger

	validator *Validator
}

type BuilderOption func(translator *Builder)

func BuilderWithLogger(logger *slog.Logger) BuilderOption {
	return func(translator *Builder) {
		translator.logger = logger
	}
}

func BuilderWithValidator(validator *Validator) BuilderOption {
	return func(translator *Builder) {
		translator.validator = validator
	}
}

func (tr *Builder) Execute(data []byte) error {
	if err := tr.validator.Execute(data); err != nil {
		return err
	}

	return nil
}

func NewBuilder(opts ...BuilderOption) *Builder {
	translator := &Builder{}
	for _, opt := range opts {
		opt(translator)
	}

	return translator
}
