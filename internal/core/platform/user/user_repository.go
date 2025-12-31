package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/uptrace/bun"
)

type Repository struct {
	logger *slog.Logger
	db     *bun.DB
}

type OptionRepository func(*Repository) error

func RepositoryWithLogger(logger *slog.Logger) OptionRepository {
	return func(r *Repository) error {
		r.logger = logger
		return nil
	}
}

func RepositoryWithDB(db *bun.DB) OptionRepository {
	return func(r *Repository) error {
		r.db = db
		return nil
	}
}

func NewRepository(opts ...OptionRepository) (*Repository, error) {
	r := &Repository{}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Repository) FindOneByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) FindAll(ctx context.Context, limit, offset int) ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
