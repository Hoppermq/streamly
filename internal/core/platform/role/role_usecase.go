package role

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/pkg/domain"
)

type UseCase struct {
	logger *slog.Logger

	roleRepo domain.RoleRepository

	generator domain.Generator
}

type UseCaseOptions func(*UseCase)

func WithLogger(logger *slog.Logger) UseCaseOptions {
	return func(u *UseCase) {
		u.logger = logger
	}
}

func WithRoleRepository(roleRepo domain.RoleRepository) UseCaseOptions {
	return func(u *UseCase) {
		u.roleRepo = roleRepo
	}
}

func WithGenerator(generator domain.Generator) UseCaseOptions {
	return func(u *UseCase) {
		u.generator = generator
	}
}

func NewUseCase(opts ...UseCaseOptions) *UseCase {
	uc := &UseCase{}

	for _, opt := range opts {
		opt(uc)
	}

	return uc
}

func (uc *UseCase) Create(ctx context.Context, input interface{}) error {
	return nil
}

func (uc *UseCase) CreateDefaults(ctx context.Context) error {
	if err := uc.Create(ctx, "default"); err != nil {
		return err
	}

	roles := uc.GenerateDefault(ctx)
	if err := uc.roleRepo.SaveAll(ctx, roles); err != nil {
		uc.logger.ErrorContext(ctx, "error during organization creation", "error", err)

		return err
	}

	return nil
}

func (uc *UseCase) GenerateDefault(ctx context.Context) []domain.TenantRole {
	uc.logger.InfoContext(ctx, "generating default tenants")

	roles := []domain.TenantRole{
		{
			Identifier:  uc.generator(),
			RoleName:    "user",
			Permissions: []string{"events:read", "dashboard:read"},
		},
		{
			Identifier:  uc.generator(),
			RoleName:    "admin",
			Permissions: []string{"events:*", "dashboard:read"},
		},
		{
			Identifier:  uc.generator(),
			RoleName:    "owner",
			Permissions: []string{"*"},
		},
	}
	return roles
}
