package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Membership struct {
	bun.BaseModel `bun:"table:tenant_members,alias:membership"`

	ID         uuid.UUID `bun:"id,type:uuid,default:uuid_generate_v4()"`
	Identifier uuid.UUID `bun:"identifier,type:uuid,notnull,unique"`

	TenantID uuid.UUID `bun:"tenant_id,type:uuid,notnull"`
	UserID   uuid.UUID `bun:"user_id,type:uuid,notnull"`

	JoinedAt time.Time `bun:"joined_at,type:timestamp"`

	CreatedAt time.Time `bun:"created_at,type:timestamp,notnull,default:now()"`
	UpdatedAt time.Time `bun:"updated_at,type:timestamp,notnull,default:now()"`
}
