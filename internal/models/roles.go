package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:tenant_roles"`

	ID         uuid.UUID `bun:"id,type:uuid,notnull,default:uuid_generate_v4()"`
	Identifier uuid.UUID `bun:"identifier,type:uuid,notnull"`

	RoleName    string         `bun:"role,notnull"`
	Permissions []string       `bun:"permissions,type:text[],notnull"`
	Metadata    map[string]any `bun:"metadata,type:jsonb"`

	CreatedAt time.Time `bun:"created_at,type:timestamp,default:now()"`
	UpdatedAt time.Time `bun:"updated_at,type:timestamp,default:now()"`
}
