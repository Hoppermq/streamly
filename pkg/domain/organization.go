package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Identifier uuid.UUID
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrganizationRepository interface {
	WithTx(tx TxContext) OrganizationRepository

	FindOneByID(ctx context.Context, identifier uuid.UUID) (*Organization, error)
	FindAll(ctx context.Context, limit, offset int) ([]Organization, error)
	Exist(ctx context.Context, identifier uuid.UUID) (bool, error)
	Create(ctx context.Context, org *Organization) error
	Update(ctx context.Context, org *Organization) error
	Delete(ctx context.Context, identifier uuid.UUID) error
}

type CreateOrganization struct {
	Name     string            `binding:"required" form:"name"`
	Metadata map[string]string `binding:"required" form:"metadata"`
}

type UpdateOrganization struct {
	Name string `binding:"required" form:"name"`
}
