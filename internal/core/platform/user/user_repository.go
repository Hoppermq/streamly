package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/internal/models"
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
	r.logger.InfoContext(ctx, "creating new user", "user", user)
	u := &models.User{
		Identifier:   user.Identifier,
		ZitadelID:    user.ZitadelID,
		Username:     user.UserName,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PrimaryEmail: user.PrimaryEmail,
		Role:         user.Role,
	}

	_, err := r.db.NewInsert().Model(u).Exec(ctx)
	if err != nil {
		r.logger.WarnContext(ctx, "failed to create a new user", "error", err)
		return err
	}

	r.logger.InfoContext(ctx, "user created successfully", "user_id", u.ID, "zitadel_id", u.ZitadelID)
	return nil

}

func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Exist(ctx context.Context, identifier uuid.UUID) (bool, error) {
	r.logger.InfoContext(ctx, "checking user existence", "user_id", identifier)

	var res bool
	_, err := r.db.NewRaw("SELECT EXISTS(SELECT 1 FROM users WHERE identifier = ? AND deleted = false);", identifier).Exec(ctx, res)
	if err != nil {
		r.logger.WarnContext(ctx, "failed to check user existence", "user_id", identifier, "error", err)
		return false, err
	}

	return res, nil
}
