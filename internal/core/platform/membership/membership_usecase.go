package membership

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/pkg/domain"
)

type UseCase struct {
	logger *slog.Logger

	repo domain.MembershipRepository

	generator  domain.Generator
	uuidParser domain.UUIDParser
}

type UseCaseOption func(*UseCase)

func UseCaseWithLogger(logger *slog.Logger) UseCaseOption {
	return func(u *UseCase) {
		u.logger = logger
	}
}

func UseCaseWithRepository(repo domain.MembershipRepository) UseCaseOption {
	return func(u *UseCase) {
		u.repo = repo
	}
}

func UseCaseWithUUIDParser(parser domain.UUIDParser) UseCaseOption {
	return func(u *UseCase) {
		u.uuidParser = parser
	}
}

func UseCaseWithGenerator(generator domain.Generator) UseCaseOption {
	return func(u *UseCase) {
		u.generator = generator
	}
}

func NewUseCase(repo any, options ...UseCaseOption) *UseCase {
	u := &UseCase{}

	for _, option := range options {
		option(u)
	}

	return u
}

func (uc *UseCase) Generate(ctx context.Context, userID, orgID string) error {
	uc.logger.Info("generating new membership", "organization", "org")
	userIdentifier, err := uc.uuidParser(userID)
	if err != nil {
		uc.logger.WarnContext(ctx, "failed to parse user identifier", "error", err)
		return err
	}

	orgIdentifier, err := uc.uuidParser(orgID)
	if err != nil {
		uc.logger.WarnContext(ctx, "failed to parse org identifier", "error", err)
		return err
	}

	identifier := uc.generator()

	membership := &domain.Membership{
		Identifier:     identifier,
		OrgIdentifier:  orgIdentifier,
		UserIdentifier: userIdentifier,
	}

	if err := uc.repo.Create(ctx, membership); err != nil {
		uc.logger.WarnContext(ctx, "failed to create new membership", "error", err)
		return err
	}

	return nil
}
