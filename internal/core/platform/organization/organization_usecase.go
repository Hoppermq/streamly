package organization

import (
	"context"
	"errors"
	"log/slog"

	"github.com/hoppermq/streamly/internal/core/platform/membership"
	"github.com/hoppermq/streamly/internal/core/platform/user"
	"github.com/hoppermq/streamly/pkg/domain"
	// should be from domain.
	"github.com/uptrace/bun"
)

type UseCase struct {
	logger *slog.Logger

	membershipUC *membership.UseCase
	userUC       *user.UseCase

	idb *bun.IDB

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

func UseCaseWithMembershipUC(membershipUC *membership.UseCase) UseCaseOption {
	return func(u *UseCase) error {
		u.membershipUC = membershipUC
		return nil
	}
}

func UseCaseWithUserUC(userUC *user.UseCase) UseCaseOption {
	return func(u *UseCase) error {
		u.userUC = userUC
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

func (uc *UseCase) Create(ctx context.Context, newOrg domain.CreateOrganization, zitadelUserID string) error {
	orgIdentifier := uc.generator()

	org := &domain.Organization{
		Identifier: orgIdentifier,
		Name:       newOrg.Name,
	}

	if err := uc.repository.Create(ctx, org); err != nil {
		uc.logger.Warn("failed to create organization", "error", err)
		return err
	}

	if err := uc.membershipUC.Generate(ctx, zitadelUserID, org.Identifier.String()); err != nil {
		uc.logger.Warn("failed to add user", "error", err)
		return err
	}

	return nil
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

func (uc *UseCase) AddUser(ctx context.Context, orgID, userID string) error {
	uc.logger.Info("adding user to organization", "orgID", orgID, "userID", userID)

	orgIdentifier, err := uc.uuidParser(orgID)
	if err != nil {
		uc.logger.Warn("failed to parse organization uuid", "error", err)
		return err
	}

	exist, err := uc.repository.Exist(ctx, orgIdentifier)
	if err != nil {
		uc.logger.Warn("failed to check organization existence", "error", err)
		return err
	}

	if !exist {
		err := errors.New("organization does not exist")

		uc.logger.Warn("failed to add user to organization", "error", err)
		return err
	}

	if err := uc.membershipUC.Generate(ctx, userID, orgID); err != nil {
		uc.logger.Warn("failed to add user to organization", "error", err)
		return err
	}

	uc.logger.Info("user added successfully", "orgID", orgID, "userID", userID)

	return nil
}
