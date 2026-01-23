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
	db     bun.IDB
}

type RepositoryOption func(*Repository)

func RepositoryWithLogger(logger *slog.Logger) RepositoryOption {
	return func(r *Repository) {
		r.logger = logger
	}
}

func RepositoryWithDB(db bun.IDB) RepositoryOption {
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

func (r *Repository) WithTx(tx domain.TxContext) domain.MembershipRepository {
	bunTx, ok := tx.(bun.IDB)
	if !ok {
		r.logger.Warn("Transaction does not implement github.com/uptrace/bun.DB")
		return r
	}

	return &Repository{
		logger: r.logger,
		db:     bunTx,
	}
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

	r.logger.InfoContext(ctx, "inserted membership", "org_id", membership.OrgIdentifier, "user_id", membership.UserIdentifier)
	return nil
}
