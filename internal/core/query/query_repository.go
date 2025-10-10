package query

import (
	"context"

	"github.com/hoppermq/streamly/pkg/domain"
)

type QueryRepository struct {
	driver domain.Driver
}

type RepositoryOption func(*QueryRepository)

func WithDriver(driver domain.Driver) RepositoryOption {
	return func(repository *QueryRepository) {
		repository.driver = driver
	}
}

func (q *QueryRepository) ExecuteQuery(ctx context.Context, req *domain.QueryAstRequest) (*domain.QueryResponse, error) {
	return nil, nil
}

func NewQueryRepository(opts ...RepositoryOption) *QueryRepository {
	qr := &QueryRepository{}
	for _, opt := range opts {
		opt(qr)
	}

	return qr
}
