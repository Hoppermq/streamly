package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:user"`
	ID            uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	Identifier    uuid.UUID `bun:"identifier,notnull,unique,type:uuid"`

	ZitadelID string `bun:"zitadel_user_id,notnull,unique"`

	FirstName    string `bun:"first_name"`
	LastName     string `bun:"last_name"`
	Username     string `bun:"username"`
	PrimaryEmail string `bun:"primary_email"`

	Role string `bun:"role"`

	Deleted bool `bun:"deleted,notnull,default:false"`

	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()"`
	DeletedAt time.Time `bun:"deleted_at,type:timestamp,default:null"`
}
