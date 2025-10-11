package query

import (
	"context"
	"database/sql"

	"github.com/hoppermq/streamly/pkg/domain"
)

type Repository struct {
	driver domain.Driver
}

type RepositoryOption func(*Repository)

func WithDriver(driver domain.Driver) RepositoryOption {
	return func(repository *Repository) {
		repository.driver = driver
	}
}

func (q *Repository) ExecuteQuery(ctx context.Context, query domain.Query, args ...domain.QueryArgs) (*domain.QueryResponse, error) {
	rows, err := q.driver.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0)

	for rows.Next() {
		values := make([]any, len(cols))
		valuePtrs := make([]any, len(values))
		for i := range cols {
			valuePtrs[i] = &values[i]
		}

		if err = rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)
		for i, col := range cols {
			rowMap[col] = values[i]
		}
		data = append(data, rowMap)
	if err := rows.Err(); err != nil {
		return nil, err
	}
}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	return &domain.QueryResponse{
		RequestID: "meta-data-user-02",
		Data:      data,
		RowCount:  len(data),
	}, nil
}

func NewQueryRepository(opts ...RepositoryOption) *Repository {
	qr := &Repository{}
	for _, opt := range opts {
		opt(qr)
	}

	return qr
}
