package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/hoppermq/streamly/internal/common"
	"github.com/hoppermq/streamly/pkg/domain"
)

type UseCase struct {
	logger *slog.Logger

	userRepo   domain.UserRepository
	authRepo   domain.AuthRepository
	generator  domain.Generator
	uuidParser domain.UUIDParser

	zitadelApi domain.Client
}

type UseCaseOption func(*UseCase) error

func UseCaseWithLogger(logger *slog.Logger) UseCaseOption {
	return func(u *UseCase) error {
		u.logger = logger
		return nil
	}
}

func UseCaseWithUserRepository(repo domain.UserRepository) UseCaseOption {
	return func(u *UseCase) error {
		u.userRepo = repo
		return nil
	}
}

func UseCaseWithAuthRepository(repo domain.AuthRepository) UseCaseOption {
	return func(u *UseCase) error {
		u.authRepo = repo
		return nil
	}
}

func UseCaseWithUUIDParser(parser domain.UUIDParser) UseCaseOption {
	return func(u *UseCase) error {
		u.uuidParser = parser
		return nil
	}
}

func UseCaseWithGenerator(generator domain.Generator) UseCaseOption {
	return func(u *UseCase) error {
		u.generator = generator
		return nil
	}
}

func UseCaseWithZitadelAPI(zitadelApi domain.Client) UseCaseOption {
	return func(u *UseCase) error {
		u.zitadelApi = zitadelApi
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
	if userInput == nil {
		// will be static error.
		err := errors.New("userInput cannot be nil")
		uc.logger.Warn("userInput is nil", "error", err)
		return err
	}
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

	var u *domain.User
	var err error
	maxRetries := 6                      // should be held in the ctx?
	retryDelay := 500 * time.Millisecond // should be held in the ctx ?

	for attempt := 1; attempt <= maxRetries; attempt++ {
		u, err = uc.zitadelApi.GetUserByUserName(ctx, event.Request.UserName)
		if err == nil {
			break
		}

		if attempt < maxRetries {
			uc.logger.Info("user not found yet, retrying...",
				"attempt", attempt,
				"username", event.Request.UserName,
				"retry_in_ms", retryDelay.Milliseconds())
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		uc.logger.Warn("failed to find user after retries",
			"email", event.Request.Email.Email,
			"attempts", maxRetries,
			"error", err.Error())
		return err
	}

	createUser := &domain.CreateUser{
		UserName:  u.UserName,
		FirstName: event.Request.Profile.FirstName,
		LastName:  event.Request.Profile.LastName,

		PrimaryEmail: event.Request.Email.Email,
		ZitadelID:    u.ZitadelID,

		Role: "user",
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

func (uc *UseCase) CheckUserExist(ctx context.Context, id string) (bool, error) {
	uc.logger.Info("checking user existence", "identifier", id)
	identifier, err := uc.uuidParser(id)
	if err != nil {
		uc.logger.Warn("failed to parse user identifier", "identifier", id, "error", err.Error())
		return false, err
	}

	return uc.userRepo.Exist(ctx, identifier)
}
