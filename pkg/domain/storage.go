package domain

import (
	"context"
)

type Storage interface {
	Save(ctx context.Context) error
	Get(ctx context.Context) error
}

type Tx interface {
	Commit() error
	Rollback() error
}

type Driver interface {
	Begin() (Tx, error)
	Close() error
}
