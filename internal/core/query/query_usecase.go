package query

import (
	"context"
	"log/slog"
	"sync"

	"github.com/hoppermq/streamly/pkg/domain"
)

// QueryUseCaseImpl represent the query use case structure.
type QueryUseCaseImpl struct {
	logger *slog.Logger
	wg     sync.WaitGroup
}

type UseCaseOption func(impl *QueryUseCaseImpl)

func UseCaseWithLogger(logger *slog.Logger) UseCaseOption {
	return func(q *QueryUseCaseImpl) {
		q.logger = logger
	}
}

func NewQueryUseCase(options ...UseCaseOption) *QueryUseCaseImpl {
	useCase := &QueryUseCaseImpl{
		wg: sync.WaitGroup{},
	}
	for _, option := range options {
		option(useCase)
	}

	return useCase
}

func (u *QueryUseCaseImpl) SyncQuery(ctx context.Context, req *domain.QueryAstRequest) (any, error) {
	return "sync query performed", nil
}

func (u *QueryUseCaseImpl) AsyncQuery(ctx context.Context, req *domain.QueryAstRequest) (any, error) {
	return "async query performed", nil
}
