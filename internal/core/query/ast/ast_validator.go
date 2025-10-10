package ast

import (
	"log/slog"
)

type Validator struct {
	logger *slog.Logger
}

type ValidatorOption func(*Validator)

func ValidatorWithLogger(logger *slog.Logger) ValidatorOption {
	return func(v *Validator) {
		v.logger = logger
	}
}

func (v *Validator) Execute(data []byte) error {
	return nil
}

func NewValidator(opts ...ValidatorOption) *Validator {
	v := &Validator{}
	for _, opt := range opts {
		opt(v)
	}

	return v
}
