package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MembershipRepository interface {
	WithTx(tx TxContext) MembershipRepository

	Create(context.Context, *Membership) error
	Update(context.Context, *Membership) error
}

type Membership struct {
	Identifier uuid.UUID

	OrgIdentifier  uuid.UUID
	UserIdentifier uuid.UUID

	RoleIdentifier uuid.UUID

	JoinedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
