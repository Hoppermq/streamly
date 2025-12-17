package domain

import (
	"context"
	"time"
)

type Organization struct {
	Identifier string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrganizationRepository interface {
	GetByID(ctx context.Context, identifier string) (*Organization, error)
	List(ctx context.Context, limit, offset int) ([]*Organization, error)
	Create(ctx context.Context, org *Organization) error
	Update(ctx context.Context, org *Organization) error
	Delete(ctx context.Context, id string) error
}

type CreateOrganization struct {
	Name     string            `form:"name"      binding:"required"`
	Metadata map[string]string `form:"metadata"     binding:"required"`
}
