package membership

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

type RepositoryOption func(*Repository)

func RepositoryWithLogger(logger *slog.Logger) RepositoryOption {
	return func(r *Repository) {
		r.logger = logger
	}
}

func RepositoryWithDB(db *bun.DB) RepositoryOption {
	return func(r *Repository) {
		r.db = db
	}
}

func NewRepository(opts ...RepositoryOption) *Repository {
	r := &Repository{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Repository) Create(ctx context.Context, membership *domain.Membership) error {
	r.logger.InfoContext(ctx, "inserting new membership", "membership_identifier", membership.Identifier)
	membershipModel := &models.Membership{
		Identifier: membership.Identifier,
		TenantID:   membership.OrgIdentifier,
		UserID:     membership.UserIdentifier,
	}

	_, err := r.db.NewInsert().Model(membershipModel).Exec(ctx)
	if err != nil {
		r.logger.Warn("failed to insert membership", "error", err)
		return err
	}

	r.logger.InfoContext(ctx, "inserted membership", membership)
	return nil

}
