package organization

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/pkg/domain"
)

type UseCase struct {
	logger *slog.Logger

	repository domain.OrganizationRepository
}

type UseCaseOption func(*UseCase) error

func UseCaseWithLogger(logger *slog.Logger) UseCaseOption {
	return func(u *UseCase) error {
		u.logger = logger
		return nil
	}
}

func UseCaseWithRepository(repository domain.OrganizationRepository) UseCaseOption {
	return func(u *UseCase) error {
		u.repository = repository
		return nil
	}
}

func NewUseCase(opts ...UseCaseOption) (*UseCase, error) {
	uc := &UseCase{}
	for _, opt := range opts {
		if err := opt(uc); err != nil {
			return nil, err
		}
	}
	return uc, nil
}

func (uc *UseCase) FindOneByID(ctx context.Context, id string) (*domain.Organization, error) {
	uc.logger.Info("finding organization by id", "id", id)
	return uc.repository.FindOneByID(ctx, id)
}

func (uc *UseCase) FindAll(ctx context.Context, limit, offset int) ([]*domain.Organization, error) {
	uc.logger.Info("finding all organizations", "limit", limit, "offset", offset)
	return uc.repository.FindAll(ctx, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, newOrg domain.CreateOrganization) error {
	orgIdentifier := uuid.New()

	org := &domain.Organization{
		Identifier: orgIdentifier,
		Name:       newOrg.Name,
	}

	return uc.repository.Create(ctx, org)
}
func (uc *UseCase) Update(ctx context.Context, id string, updateOrg domain.UpdateOrganization) error {
	org, err := uc.repository.FindOneByID(ctx, id)
	if err != nil {
		return err
	}

	org.Name = updateOrg.Name

	return uc.repository.Update(ctx, org)
}
