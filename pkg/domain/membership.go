package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MembershipRepository interface {
	WithTx(tx TxContext) MembershipRepository

	Create(ctx context.Context, membership *Membership) error
}

type Membership struct {
	Identifier uuid.UUID

	OrgIdentifier  uuid.UUID
	UserIdentifier uuid.UUID

	JoinedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
