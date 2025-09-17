package domain

import (
	"context"
	"database/sql"
)

type Storage interface {
	Save(ctx context.Context) error
	Get(ctx context.Context) error
}

type Stmt interface {
	ExecContext(ctx context.Context, args ...interface{}) error
	Close() error
}

type Tx interface {
	Commit() error
	Rollback() error
	PrepareContext(ctx context.Context, query string) (Stmt, error)
}

type Driver interface {
	Begin() (Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) // TXOption should be an interface.
	Close() error
}
