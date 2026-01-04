package user

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/internal/common"
	"github.com/hoppermq/streamly/pkg/domain"
)

type UseCase struct {
	logger *slog.Logger

	userRepo   domain.UserRepository
	authRepo   domain.AuthRepository
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

func WithUserRepository(repo domain.UserRepository) UseCaseOption {
	return func(u *UseCase) error {
		u.userRepo = repo
		return nil
	}
}

func WithAuthRepository(repo domain.AuthRepository) UseCaseOption {
	return func(u *UseCase) error {
		u.authRepo = repo
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

	return uc.userRepo.FindOneByID(ctx, identifier)
}

func (uc *UseCase) FindAll(ctx context.Context, limit, offset int) ([]domain.User, error) {
	uc.logger.Info("finding users", "limit", limit, "offset", offset)
	return uc.userRepo.FindAll(ctx, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, userInput *domain.CreateUser) error {
	uc.logger.Info("creating new user")
	userIdentifier := uc.generator()

	user := &domain.User{
		Identifier: userIdentifier,
		ZitadelID:  userInput.ZitadelID,

		UserName:     userInput.UserName,
		FirstName:    userInput.FirstName,
		LastName:     userInput.LastName,
		PrimaryEmail: userInput.PrimaryEmail,

		Role: userInput.Role,
	}

	return uc.userRepo.Create(ctx, user)
}

func (uc *UseCase) CreateFromEvent(ctx context.Context, event *domain.ZitadelEventUserCreated) error {
	uc.logger.Info("creating new user from event")
	u, err := uc.authRepo.FindUserByUsername(ctx, event.Request.UserName)
	if err != nil {
		uc.logger.Warn("failed to find user by email", "email", event.Request.Email.Email, "error", err.Error())
		return err
	}

	createUser := &domain.CreateUser{
		UserName:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,

		PrimaryEmail: u.PrimaryEmail,
		ZitadelID:    u.ZitadelID,

		Role: u.Role,
	}

	return uc.Create(ctx, createUser)
}

func (uc *UseCase) Update(ctx context.Context, id string, updateUserInput *domain.UpdateUser) error {
	uc.logger.Info("updating user by id", "identifier", id)
	identifier, err := uc.uuidParser(id)
	if err != nil {
		uc.logger.Warn("failed to parse user identifier", "identifier", id, "error", err.Error())
		return err
	}

	existingUser, err := uc.userRepo.FindOneByID(ctx, identifier)
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

	return uc.userRepo.Update(ctx, existingUser)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	uc.logger.Info("deleting user by id", "identifier", id)
	identifier, err := uc.uuidParser(id)
	if err != nil {
		uc.logger.Warn("failed to parse user identifier", "identifier", id, "error", err.Error())
		return err
	}

	return uc.userRepo.Delete(ctx, identifier)
}
