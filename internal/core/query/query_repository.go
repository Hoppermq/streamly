package query

import (
	"context"

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

func (q *Repository) ExecuteQuery(ctx context.Context, req *domain.QueryAstRequest) (*domain.QueryResponse, error) {
	return &domain.QueryResponse{
		RequestID: req.RequestID,
		Data:      []map[string]any{},
		RowCount:  0,
	}, nil
}

func NewQueryRepository(opts ...RepositoryOption) *Repository {
	qr := &Repository{}
	for _, opt := range opts {
		opt(qr)
	}

	return qr
}
