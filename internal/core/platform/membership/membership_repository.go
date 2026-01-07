package membership

import (
	"context"
	"errors"
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
		RoleID:     membership.RoleIdentifier,
	}

	_, err := r.db.NewInsert().Model(membershipModel).Exec(ctx)
	if err != nil {
		r.logger.Warn("failed to insert membership", "error", err)
		return err
	}

	r.logger.InfoContext(ctx, "inserted membership", membership)
	return nil

}

func (r *Repository) Update(ctx context.Context, membership *domain.Membership) error {
	r.logger.InfoContext(ctx, "updating membership", "membership_id", membership.Identifier)
	model := &models.Membership{
		Identifier: membership.Identifier,
		TenantID:   membership.OrgIdentifier,
		UserID:     membership.UserIdentifier,
		JoinedAt:   membership.JoinedAt,
	}

	res, err := r.db.NewUpdate().
		Model(model).
		Where("identifier = ?", membership.Identifier).
		Where("deleted = ?", false).
		Exec(ctx)

	if err != nil {
		r.logger.Warn("failed to update membership", "error", err)
		return err
	}

	if res, err := res.RowsAffected(); err != nil || res == 0 {
		if err != nil {
			r.logger.Warn("failed to update membership", "error", err)
			return err
		}
		return errors.New("failed to update membership")
	}

	return nil
}
