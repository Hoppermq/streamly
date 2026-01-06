package organization

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/internal/models"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/uptrace/bun"
)

type Repository struct {
	logger *slog.Logger
	db     bun.IDB
}

type OptionRepository func(*Repository) error

func RepositoryWithLogger(logger *slog.Logger) OptionRepository {
	return func(organizationRepo *Repository) error {
		organizationRepo.logger = logger
		return nil
	}
}

func RepositoryWithDB(db *bun.DB) OptionRepository {
	return func(organizationRepo *Repository) error {
		organizationRepo.db = db
		return nil
	}
}

func NewRepository(opts ...OptionRepository) (*Repository, error) {
	org := &Repository{}

	for _, opt := range opts {
		if err := opt(org); err != nil {
			return nil, err
		}
	}

	return org, nil
}

func (organizationRepo *Repository) WithTx(tx interface{}) domain.OrganizationRepository {
	bunTx, ok := tx.(bun.IDB)
	if !ok {
		organizationRepo.logger.Warn("Transaction does not implement github.com/uptrace/bun.DB")
		return organizationRepo
	}

	return &Repository{
		logger: organizationRepo.logger,
		db:     bunTx,
	}
}

func (organizationRepo *Repository) FindOneByID(
	ctx context.Context,
	identifier uuid.UUID,
) (*domain.Organization, error) {
	organizationRepo.logger.InfoContext(ctx, "getting organization from id", "id", identifier)

	org := &models.Organization{}

	if err := organizationRepo.db.NewSelect().Model(org).Where("identifier = ?", identifier).Where("deleted = ?", false).Scan(ctx); err != nil {
		organizationRepo.logger.WarnContext(ctx, "failed to select org", "identifier", identifier, "error", err)
		return nil, err
	}

	organizationRepo.logger.Info("organization", "data", org)
	if org.Identifier == uuid.Nil {
		return nil, errors.New("organization not found") // TODO: return a custom error type for not found.
	}

	organizationRepo.logger.InfoContext(
		ctx,
		"organization found",
		"identifier",
		org.Identifier,
		"name",
		org.Name,
	)

	res := domain.Organization{
		Identifier: org.Identifier,
		Name:       org.Name,
		CreatedAt:  org.CreatedAt,
		UpdatedAt:  org.UpdatedAt,
	}

	return &res, nil
}

func (organizationRepo *Repository) FindAll(
	ctx context.Context,
	limit, offset int,
) ([]domain.Organization, error) {
	var orgs []models.Organization
	if err := organizationRepo.db.NewSelect().Model(&orgs).Where("deleted = ?", false).Limit(limit).Offset(offset).Scan(ctx); err != nil {
		organizationRepo.logger.Warn("failed to query tenants", "error", err)
		return nil, err
	}

	organizations := make([]domain.Organization, len(orgs))
	for i, org := range orgs {
		organizations[i] = domain.Organization{
			Identifier: org.Identifier,
			Name:       org.Name,
			CreatedAt:  org.CreatedAt,
			UpdatedAt:  org.UpdatedAt,
		}
	}

	return organizations, nil
}

func (organizationRepo *Repository) Create(
	ctx context.Context,
	org *domain.Organization,
) error {
	organizationRepo.logger.InfoContext(ctx, "inserting new org", "org_identifier", org.Identifier)
	model := &models.Organization{
		Identifier: org.Identifier,
		Name:       org.Name,
	}
	res, err := organizationRepo.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		organizationRepo.logger.WarnContext(ctx, "failed to insert new org", "error", err)
		return err
	}

	organizationRepo.logger.InfoContext(ctx, "organization insertion successed", "organization", res)
	return nil
}

func (organizationRepo *Repository) Update(
	ctx context.Context,
	org *domain.Organization,
) error {
	organizationRepo.logger.InfoContext(ctx, "updating org", "org_identifier", org.Identifier)
	model := &models.Organization{
		Identifier: org.Identifier,
		Name:       org.Name,
		UpdatedAt:  org.UpdatedAt,
	}

	res, err := organizationRepo.
		db.NewUpdate().
		Model(model).
		Column("name", "updated_at").
		Where("identifier = ?", org.Identifier).
		Where("deleted = ?", false).
		Exec(ctx)

	if err != nil {
		organizationRepo.logger.WarnContext(ctx, "failed to update org", "error", err)
		return err
	}

	if res, err := res.RowsAffected(); res == 0 || err != nil {
		if err != nil {
			organizationRepo.logger.WarnContext(ctx, "failed to get rows affected", "error", err)
		}
		return errors.New("failed to update org")
	}

	return nil
}

func (organizationRepo *Repository) Delete(
	ctx context.Context,
	identifier uuid.UUID,
) error {
	organizationRepo.logger.InfoContext(ctx, "deleting org", "org_id", identifier)
	org := &models.Organization{
		Identifier: identifier,
	}

	res, err := organizationRepo.
		db.
		NewUpdate().
		Model(org).
		Where("identifier = ?", identifier).
		Where("deleted = ?", false).
		Set("deleted_at = ?", time.Now()).
		Set("deleted = ?", true).
		Exec(ctx)

	if err != nil {
		organizationRepo.logger.WarnContext(ctx, "failed to delete org", "error", err)
		return err
	}

	if res, err := res.RowsAffected(); res == 0 || err != nil {
		if err != nil {
			organizationRepo.logger.WarnContext(ctx, "failed to get rows affected", "error", err)
		}
		return errors.New("failed to delete org")
	}

	return nil
}

func (organizationRepo *Repository) Exist(ctx context.Context, identifier uuid.UUID) (bool, error) {
	var res bool
	err := organizationRepo.db.NewRaw(
		"SELECT EXISTS(SELECT 1 FROM tenants WHERE identifier = ? AND deleted = false);",
		identifier,
	).Scan(ctx, &res)
	if err != nil {
		organizationRepo.logger.WarnContext(ctx, "failed to query tenants", "error", err)
		return false, err
	}

	return res, nil
}
