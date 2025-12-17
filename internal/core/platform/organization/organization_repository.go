package organization

import (
	"context"
	"errors"
	"log/slog"

	"github.com/hoppermq/streamly/internal/models"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/uptrace/bun"
)

type OrganizationRepository struct {
	logger *slog.Logger
	db     *bun.DB
}

type OptionRepository func(*OrganizationRepository) error

func RepositoryWithLogger(logger *slog.Logger) OptionRepository {
	return func(organizationRepo *OrganizationRepository) error {
		organizationRepo.logger = logger
		return nil
	}
}

func NewRepository(opts ...OptionRepository) (*OrganizationRepository, error) {
	org := &OrganizationRepository{}

	for _, opt := range opts {
		if err := opt(org); err != nil {
			return nil, err
		}
	}

	return org, nil
}

func (organizationRepo *OrganizationRepository) GetByID(
	ctx context.Context,
	identifier string,
) (*domain.Organization, error) {
	organizationRepo.logger.InfoContext(ctx, "getting organization from id", "id", identifier)

	org := &models.Organization{}

	if err := organizationRepo.db.NewSelect().Model(&org).Where("identifier = ?", identifier).Scan(ctx); err != nil {
		organizationRepo.logger.WarnContext(ctx, "failed to select org", "identifier", identifier, "error", err)
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
		Metada:     org.Metadata,
		CreatedAt:  org.CreatedAt,
		UpdatedAt:  org.UpdatedAt,
	}

	return &res, nil
}

func (organizationRepo *OrganizationRepository) List(
	ctx context.Context,
	limit, offset int,
) ([]*domain.Organization, error) {
	return nil, nil
}

func (organizationRepo *OrganizationRepository) Create(
	ctx context.Context,
	org *domain.Organization,
) error {
	organizationRepo.logger.InfoContext(ctx, "inserting new org", "org_identifier", org.Identifier)
	model := &models.Organization{
		Identifier: org.Identifier,
		Name:       org.Name,
		Metadata:   org.Metada,
	}
	res, err := organizationRepo.db.NewInsert().Model(&model).Exec(ctx)
	if err != nil {
		organizationRepo.logger.WarnContext(ctx, "failed to insert new org", "error", err)
		return err
	}

	organizationRepo.logger.InfoContext(ctx, "organization insertion successed", "organization", res)
	return nil
}

func (organizationRepo *OrganizationRepository) Update(
	ctx context.Context,
	org *domain.Organization,
) error {
	organizationRepo.logger.InfoContext(ctx, "updating org", "org_identifier", org.Identifier)
	model := &models.Organization{
		Identifier: org.Identifier,
		Name:       org.Name,
		Metadata:   org.Metada,
		UpdatedAt:  org.UpdatedAt,
	}

	res, err := organizationRepo.
		db.NewUpdate().
		Model(&model).
		Column("name", "metadata", "updated_at").
		Where("identifier = ?", org.Identifier).
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

func (organizationRepo *OrganizationRepository) Delete(
	ctx context.Context,
	identifier string,
) error {
	return nil
}
