package organization

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/pkg/domain"
)

type UseCase struct {
	logger *slog.Logger

	repository domain.OrganizationRepository
	generator  domain.Generator
	uuidParser domain.UUIDParser
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

func UseCaseWithGenerator(generator domain.Generator) UseCaseOption {
	return func(u *UseCase) error {
		u.generator = generator
		return nil
	}
}

func UseCaseWithUUIDParser(uuidParser domain.UUIDParser) UseCaseOption {
	return func(u *UseCase) error {
		u.uuidParser = uuidParser
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
	identifier, err := uc.uuidParser(id)
	if err != nil {
		uc.logger.Warn("failed to parse uuid", "error", err)
		return nil, err
	}
	return uc.repository.FindOneByID(ctx, identifier)
}

func (uc *UseCase) FindAll(ctx context.Context, limit, offset int) ([]domain.Organization, error) {
	uc.logger.Info("finding all organizations", "limit", limit, "offset", offset)
	return uc.repository.FindAll(ctx, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, newOrg domain.CreateOrganization) error {
	orgIdentifier := uc.generator()

	org := &domain.Organization{
		Identifier: orgIdentifier,
		Name:       newOrg.Name,
	}

	return uc.repository.Create(ctx, org)
}

func (uc *UseCase) Update(ctx context.Context, id string, updateOrg domain.UpdateOrganization) error {
	identifier, err := uc.uuidParser(id)
	if err != nil {
		uc.logger.Warn("failed to parse uuid", "error", err)
		return err
	}

	org, err := uc.repository.FindOneByID(ctx, identifier)
	if err != nil {
		return err
	}

	org.Name = updateOrg.Name

	return uc.repository.Update(ctx, org)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	uc.logger.Info("deleting organization", "id", id)
	identifier, err := uc.uuidParser(id)
	if err != nil {
		uc.logger.Warn("failed to parse uuid", "error", err)
		return err
	}
	return uc.repository.Delete(ctx, identifier)
}
