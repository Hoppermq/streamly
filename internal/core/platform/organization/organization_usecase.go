package organization

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/internal/core/platform/membership"
	"github.com/hoppermq/streamly/internal/core/platform/user"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/errors"
)

type UseCase struct {
	logger *slog.Logger

	membershipUC *membership.UseCase
	userUC       *user.UseCase

	repository domain.OrganizationRepository
	generator  domain.Generator
	uuidParser domain.UUIDParser

	uow domain.UnitOfWorkFactory
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

func UseCaseWithUOW(uow domain.UnitOfWorkFactory) UseCaseOption {
	return func(u *UseCase) error {
		u.uow = uow
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
	uow, err := uc.uow.NewUnitOfWork(ctx)
	if err != nil {
		uc.logger.WarnContext(ctx, "failed to create unit of work", "error", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = uow.Rollback()
			uc.logger.ErrorContext(ctx, "panic during organization creation", "panic", p)
			// should be gracefully here.
			panic(p)
		}
	}()

	orgIdentifier := uc.generator()
	org := &domain.Organization{
		Identifier: orgIdentifier,
		Name:       newOrg.Name,
	}

	if err = uow.Organization().Create(ctx, org); err != nil {
		uc.logger.WarnContext(ctx, "failed to create organization", "error", err)
		_ = uow.Rollback()
		return err
	}

	userIdentifier, err := uow.User().GetUserIDFromZitadelID(ctx, zitadelUserID)
	if err != nil {
		uc.logger.WarnContext(ctx, "failed to get user identifier", "error", err)
		_ = uow.Rollback()
		return err
	}

	m := &domain.Membership{
		Identifier:     uc.generator(),
		OrgIdentifier:  orgIdentifier,
		UserIdentifier: userIdentifier,
	}

	if err = uow.Membership().Create(ctx, m); err != nil {
		uc.logger.WarnContext(ctx, "failed to create membership", "error", err)
		_ = uow.Rollback()
		return err
	}

	if err := uow.Commit(); err != nil {
		uc.logger.WarnContext(ctx, "failed to commit transaction", "error", err)
		return err
	}

	uc.logger.InfoContext(ctx, "organization created successfully",
		"org_id", orgIdentifier,
		"user_id", zitadelUserID,
	)

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
		err := errors.ErrOrganizationNotFound

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
