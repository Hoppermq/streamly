package auth

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/internal/models"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/uptrace/bun"
)

type Repository struct {
	logger *slog.Logger
	db     *bun.DB
}

type RepositoryOption func(*Repository) error

func RepositoryWithLogger(logger *slog.Logger) RepositoryOption {
	return func(r *Repository) error {
		r.logger = logger
		return nil
	}
}

func RepositoryWithDB(db *bun.DB) RepositoryOption {
	return func(r *Repository) error {
		r.db = db
		return nil
	}
}

func NewRepository(opts ...RepositoryOption) (*Repository, error) {
	repository := &Repository{}
	for _, opt := range opts {
		if err := opt(repository); err != nil {
			return nil, err
		}
	}

	return repository, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	u := &models.UserAuth{}
	if err := r.db.NewSelect().Model(u).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	foundUser := &domain.User{ZitadelID: u.ID}
	return foundUser, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &models.UserAuth{}
	if err := r.db.NewSelect().Model(u).Where("email = ?", email).Scan(ctx); err != nil {
		r.logger.WarnContext(ctx, "failed to find user by email", "email", email)
		return nil, err
	}

	foundUser := &domain.User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		PrimaryEmail: u.Email,
		ZitadelID:    u.ID,
	}
	return foundUser, nil
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	u := &models.UserProjection{}

	if err := r.db.NewSelect().Model(u).Where("username = ?", username).Scan(ctx); err != nil {
		r.logger.WarnContext(ctx, "failed to find user by username", "username", username)
		return nil, err
	}

	foundUser := &domain.User{
		ZitadelID: u.ID,
		UserName:  u.UserName,
	}

	return foundUser, nil
}
