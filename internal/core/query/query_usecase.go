package query

import (
	"context"
	"encoding/json"
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

	data, err := json.Marshal(req)
	if err != nil {
		uc.logger.Warn("failed to marshal query request", "error", err.Error())
		return nil, err
	}

	err = uc.astBuilder.Execute(data)
	if err != nil {
		uc.logger.Warn("error while building the query ast", "error", err.Error())
		return nil, err
	}

	if uc.repository != nil {
		return uc.repository.ExecuteQuery(ctx, req)
	}

	return &domain.QueryResponse{
		RequestID: req.RequestID,
		Data:      []map[string]any{},
		RowCount:  0,
	}, nil
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
