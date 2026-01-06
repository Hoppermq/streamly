package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Membership struct {
	bun.BaseModel `bun:"table:tenants"`

	ID         uuid.UUID `bun:"id,type:uuid,default:uuid_generate_v4()"`
	Identifier uuid.UUID `bun:"identifier,type:uuid,notnull,unique"`

	TenantID uuid.UUID `bun:"tenant_id,type:uuid,notnull"`
	UserID   uuid.UUID `bun:"user_id,type:uuid,notnull"`

	JoinedAt time.Time `bun:"joinedAt,type:timestamp"`

	CreatedAt time.Time `bun:"createdAt,type:timestamp,notnull,default:now()"`
	UpdatedAt time.Time `bun:"updatedAt,type:timestamp,notnull,default:now()"`
}
