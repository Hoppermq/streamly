package clickhouse

import (
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hoppermq/streamly/cmd/config"
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

type DriverOption func(options *clickhouse.Options)

func WithConfig(clickhouseConfig *config.IngestionConfig) DriverOption {
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

func OpenConn(opts ...DriverOption) domain.Driver {
	options := &clickhouse.Options{}
	for _, opt := range opts {
		opt(options)
	}

	fmt.Println("HELLO ?", options.Addr)
	db := clickhouse.OpenDB(options)
	return &ClickHouseDriver{db: db}
}
