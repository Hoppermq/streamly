package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/pkg/domain"
)

func (z *Zitadel) GetOrganizationByID(ctx context.Context, organizationId uuid.UUID) (*domain.Organization, error) {
	// TODO implement me
	panic("implement me")
}
