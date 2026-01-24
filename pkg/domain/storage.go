package domain

import (
	"context"
	"database/sql"
	"time"
)

type TokenCacheKey string
type Cache[T any] interface {
	Get(token string) (v T, ok bool)
	Set(token string, value T, ttl time.Duration)
}

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
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) // TXOption should be an interface.
	Close() error
	Query(ctx context.Context, query Query, args ...QueryArgs) (*sql.Rows, error)
	QueryContext(ctx context.Context, query Query, args ...QueryArgs) (*sql.Rows, error)
}

type UnitOfWork interface {
	Commit() error
	Rollback() error

	Organization() OrganizationRepository
	User() UserRepository
	Membership() MembershipRepository
}

type TxContext interface{}

// UnitOfWorkFactory creates new UnitOfWork instances (transactions).
type UnitOfWorkFactory interface {
	NewUnitOfWork(ctx context.Context) (UnitOfWork, error)
}
