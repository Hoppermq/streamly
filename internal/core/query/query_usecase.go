package query

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/pkg/domain"
)

type QueryUseCaseImpl struct {
	logger     *slog.Logger
	repository domain.QueryRepository
}

type UseCaseOption func(*QueryUseCaseImpl)

func NewQueryUseCase(opts ...UseCaseOption) *QueryUseCaseImpl {
	uc := &QueryUseCaseImpl{}
	for _, opt := range opts {
		opt(uc)
	}
	return uc
}

func WithUseCaseLogger(logger *slog.Logger) UseCaseOption {
	return func(uc *QueryUseCaseImpl) {
		uc.logger = logger
	}
}

func WithRepository(repo domain.QueryRepository) UseCaseOption {
	return func(uc *QueryUseCaseImpl) {
		uc.repository = repo
	}
}

func (uc *QueryUseCaseImpl) SyncQuery(ctx context.Context, req *domain.QueryAstRequest) (*domain.QueryResponse, error) {
	uc.applyDefaults(req)

	if uc.repository != nil {
		return uc.repository.ExecuteQuery(ctx, req)
	}

	return &domain.QueryResponse{
		RequestID: req.RequestID,
		Data:      []map[string]any{},
		RowCount:  0,
	}, nil
}

func (uc *QueryUseCaseImpl) applyDefaults(req *domain.QueryAstRequest) {
	if req.Limit == nil {
		defaultLimit := 1000
		req.Limit = &defaultLimit
	}

	if req.Offset == nil {
		defaultOffset := 0
		req.Offset = &defaultOffset
	}

	for i := range req.OrderBy {
		if req.OrderBy[i].Direction == "" {
			req.OrderBy[i].Direction = "DESC"
		}
	}
}
