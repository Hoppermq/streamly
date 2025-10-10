package query

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/internal/core/query/ast"
	"github.com/hoppermq/streamly/pkg/domain"
)

type UseCaseImpl struct {
	logger     *slog.Logger
	repository domain.QueryRepository

	astBuilder *ast.Builder
}

type UseCaseOption func(*UseCaseImpl)

func NewQueryUseCase(opts ...UseCaseOption) *UseCaseImpl {
	uc := &UseCaseImpl{}
	for _, opt := range opts {
		opt(uc)
	}
	return uc
}

func WithUseCaseLogger(logger *slog.Logger) UseCaseOption {
	return func(uc *UseCaseImpl) {
		uc.logger = logger
	}
}

func WithRepository(repo domain.QueryRepository) UseCaseOption {
	return func(uc *UseCaseImpl) {
		uc.repository = repo
	}
}

func WithAstBuilder(astBuilder *ast.Builder) UseCaseOption {
	return func(uc *UseCaseImpl) {
		uc.astBuilder = astBuilder
	}
}

func (uc *UseCaseImpl) SyncQuery(ctx context.Context, req *domain.QueryAstRequest) (*domain.QueryResponse, error) {
	applyDefaults(req)

	err := uc.astBuilder.Execute(req)
	if err != nil {
		uc.logger.Warn("error while building the query ast", "error", err.Error())
		return nil, err
	}

	return uc.repository.ExecuteQuery(ctx, req)
}

func applyDefaults(req *domain.QueryAstRequest) {
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
