package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/internal/models"
)

type TenantRole struct {
	Identifier uuid.UUID

	RoleName    string
	Permissions []string
	Metadata    map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTenantRoleInput struct {
	Identifier  uuid.UUID
	RoleName    string
	Permissions []string
	Metadata    map[string]any
}

type UpdateTenantRoleInput struct{}

type RoleRepository interface {
	WithTx(tx TxContext) RoleRepository

	Create(context.Context, CreateTenantRoleInput) error
	Save(context.Context, models.Role) error
	SaveAll(context.Context, []TenantRole) error
	Update(context.Context, UpdateTenantRoleInput) error
	FindByID(ctx context.Context, id uuid.UUID) (*TenantRole, error)
	FindAll(ctx context.Context, limit, offset int) ([]TenantRole, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
