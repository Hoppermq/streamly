package organization

import (
	"context"
	"log/slog"

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

func (uc *UseCase) Create(ctx context.Context, newOrg domain.CreateOrganization) error {
	org := &domain.Organization{
		Identifier: "new-identifier-generated",
		Name:       newOrg.Name,
	}

	return uc.repository.Create(ctx, org)
}
