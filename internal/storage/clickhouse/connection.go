package clickhouse

import (
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hoppermq/streamly/pkg/domain"
)

// ClickHouseDriver adapts *sql.DB to domain.Driver interface
type ClickHouseDriver struct {
	db *sql.DB
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

func (d *ClickHouseDriver) Begin() (domain.Tx, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	return &txAdapter{tx: tx}, nil
}

func (d *ClickHouseDriver) Close() error {
	return d.db.Close()
}

// should have a connection type for the toml config.
type DriverOption func(options *clickhouse.Options)

func WithAddr(addrs ...string) DriverOption {
	return func(options *clickhouse.Options) {
		for _, addr := range addrs {
			options.Addr = append(options.Addr, addr)
		}
	}
}

func WithUser(user string) DriverOption {
	return func(options *clickhouse.Options) {
		options.Auth.Username = user
	}
}

func WithPassword(pw string) DriverOption {
	return func(options *clickhouse.Options) {
		options.Auth.Password = pw
	}
}

func WithDatabase(database string) DriverOption {
	return func(options *clickhouse.Options) {
		options.Auth.Database = database
	}
}

func OpenConn(opts ...DriverOption) domain.Driver {
	options := &clickhouse.Options{}
	for _, opt := range opts {
		opt(options)
	}

	db := clickhouse.OpenDB(options)
	return &ClickHouseDriver{db: db}
}
