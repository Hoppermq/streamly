package user

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/internal/common"
	"github.com/hoppermq/streamly/pkg/domain"
)

type UseCase struct {
	logger *slog.Logger

	repo       domain.UserRepository
	generator  domain.Generator
	uuidParser domain.UUIDParser
}

type UseCaseOption func(*UseCase) error

func WithLogger(logger *slog.Logger) UseCaseOption {
	return func(u *UseCase) error {
		u.logger = logger
		return nil
	}
}

func WithRepository(repo domain.UserRepository) UseCaseOption {
	return func(u *UseCase) error {
		u.repo = repo
		return nil
	}
}

func WithUUIDParser(parser domain.UUIDParser) UseCaseOption {
	return func(u *UseCase) error {
		u.uuidParser = parser
		return nil
	}
}

func WithGenerator(generator domain.Generator) UseCaseOption {
	return func(u *UseCase) error {
		u.generator = generator
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

func (uc *UseCase) FindOne(ctx context.Context, id string) (*domain.User, error) {
	uc.logger.Info("finding user by id", "identifier", id)
	identifier, err := uc.uuidParser(id)

	if err != nil {
		uc.logger.Warn("failed to parse user identifier", "identifier", id, "error", err.Error())
		return nil, err
	}

	return uc.repo.FindOneByID(ctx, identifier)
}

func (uc *UseCase) FindAll(ctx context.Context, limit, offset int) ([]domain.User, error) {
	uc.logger.Info("finding users", "limit", limit, "offset", offset)
	return uc.repo.FindAll(ctx, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, userInput *domain.CreateUser) error {
	uc.logger.Info("creating new user")
	userIdentifier := uc.generator()
	zitadelIdentifier := uc.generator()

	user := &domain.User{
		Identifier: userIdentifier,
		ZitadelID:  zitadelIdentifier,

		UserName:     userInput.UserName,
		FirstName:    userInput.FirstName,
		LastName:     userInput.LastName,
		PrimaryEmail: userInput.PrimaryEmail,

		Role: userInput.Role,
	}

	return uc.repo.Create(ctx, user)
}

func (uc *UseCase) Update(ctx context.Context, id string, updateUserInput *domain.UpdateUser) error {
	uc.logger.Info("updating user by id", "identifier", id)
	identifier, err := uc.uuidParser(id)
	if err != nil {
		uc.logger.Warn("failed to parse user identifier", "identifier", id, "error", err.Error())
		return err
	}

	existingUser, err := uc.repo.FindOneByID(ctx, identifier)
	if err != nil {
		uc.logger.Warn("failed to find user by id", "identifier", id, "error", err.Error())
		return err
	}

	var updateFields []string

	updateFields = common.ApplyStringUpdate(
		&existingUser.UserName,
		&updateUserInput.UserName,
		"first_name",
		updateFields,
	)
	updateFields = common.ApplyStringUpdate(
		&existingUser.FirstName,
		&updateUserInput.FirstName,
		"first_name",
		updateFields,
	)
	updateFields = common.ApplyStringUpdate(
		&existingUser.LastName,
		&updateUserInput.LastName,
		"last_name",
		updateFields,
	)
	updateFields = common.ApplyStringUpdate(
		&existingUser.PrimaryEmail,
		&updateUserInput.PrimaryEmail,
		"primary_email",
		updateFields,
	)

	if len(updateFields) == 0 {
		uc.logger.Warn("no chages detected for user", "identifier", id)
		return nil
	}

	uc.logger.Info("updating user fields",
		"identifier", id,
		"updated_fields", updateFields)

	return uc.repo.Update(ctx, existingUser)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	uc.logger.Info("deleting user by id", "identifier", id)
	identifier, err := uc.uuidParser(id)
	if err != nil {
		uc.logger.Warn("failed to parse user identifier", "identifier", id, "error", err.Error())
		return err
	}

	return uc.repo.Delete(ctx, identifier)
}
