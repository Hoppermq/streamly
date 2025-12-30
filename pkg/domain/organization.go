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
	FindOneByID(ctx context.Context, identifier string) (*Organization, error)
	FindAll(ctx context.Context, limit, offset int) ([]Organization, error)
	Create(ctx context.Context, org *Organization) error
	Update(ctx context.Context, org *Organization) error
	Delete(ctx context.Context, org *Organization) error
}

type CreateOrganization struct {
	Name     string            `form:"name"      binding:"required"`
	Metadata map[string]string `form:"metadata"     binding:"required"`
}

type UpdateOrganization struct {
	Name string `form:"name" binding:"required"`
}
