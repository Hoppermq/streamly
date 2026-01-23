package clickhouse

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/pkg/domain"
)

// Driver adapts *sql.DB to domain.Driver interface.
type Driver struct {
	db *sql.DB
}

type stmtAdapter struct {
	stmt *sql.Stmt
}

func (s *stmtAdapter) ExecContext(ctx context.Context, args ...interface{}) error {
	_, err := s.stmt.ExecContext(ctx, args...)
	return err
}

func (s *stmtAdapter) Close() error {
	return s.stmt.Close()
}

type txAdapter struct {
	tx *sql.Tx
}

func (t *txAdapter) Commit() error {
	return t.tx.Commit()
}

func (t *txAdapter) Rollback() error {
	return t.tx.Rollback()
}

func (t *txAdapter) PrepareContext(ctx context.Context, query string) (domain.Stmt, error) {
	stmt, err := t.tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &stmtAdapter{stmt: stmt}, nil
}

func (d *Driver) BeginTx(ctx context.Context, opts *sql.TxOptions) (domain.Tx, error) {
	tx, err := d.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &txAdapter{tx: tx}, nil
}

func (d *Driver) Close() error {
	return d.db.Close()
}

func (d *Driver) DB() *sql.DB {
	return d.db
}

func (d *Driver) Query(ctx context.Context, query domain.Query, args ...domain.QueryArgs) (*sql.Rows, error) {
	compliantArgs := make([]any, len(args))
	for i, a := range args {
		compliantArgs[i] = a
	}
	fmt.Println(compliantArgs)
	return d.db.QueryContext(ctx, string(query), compliantArgs...)
}

func (d *Driver) QueryContext(ctx context.Context, query domain.Query, args ...domain.QueryArgs) (*sql.Rows, error) {
	compliantArgs := make([]any, len(args))
	for i, a := range args {
		compliantArgs[i] = a
	}

	fmt.Println(compliantArgs)
	return d.db.QueryContext(ctx, string(query), compliantArgs...)
}

type DriverOption func(options *clickhouse.Options)

// WithIngestionConfig TODO: extract since it's pure domain logic.
func WithIngestionConfig(clickhouseConfig *config.IngestionConfig) DriverOption {
	return func(options *clickhouse.Options) {
		options.Addr = []string{
			clickhouseConfig.Ingestor.Storage.Clickhouse.Address +
				":" + clickhouseConfig.Ingestor.Storage.Clickhouse.Port,
		}
		options.Auth.Database = clickhouseConfig.Ingestor.Storage.Clickhouse.Database
		options.Auth.Username = clickhouseConfig.Ingestor.Storage.Clickhouse.UserName
		options.Auth.Password = clickhouseConfig.Ingestor.Storage.Clickhouse.Password
	}
}

// WithQueryConfig TODO: extract since it's pure domain logic.
func WithQueryConfig(clickhouseConfig *config.QueryConfig) DriverOption {
	return func(options *clickhouse.Options) {
		options.Addr = []string{
			clickhouseConfig.Query.Storage.Clickhouse.Address + ":" + clickhouseConfig.Query.Storage.Clickhouse.Port,
		}
		options.Auth.Database = clickhouseConfig.Query.Storage.Clickhouse.Database
		options.Auth.Username = clickhouseConfig.Query.Storage.Clickhouse.UserName
		options.Auth.Password = clickhouseConfig.Query.Storage.Clickhouse.Password
	}
}

func OpenConn(opts ...DriverOption) domain.Driver {
	options := &clickhouse.Options{}
	for _, opt := range opts {
		opt(options)
	}

	db := clickhouse.OpenDB(options)
	return &Driver{db: db}
}
