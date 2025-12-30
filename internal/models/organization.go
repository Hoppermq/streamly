package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Organization struct {
	bun.BaseModel `bun:"table:tenants,alias:org"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,scanonly"`
	Identifier    uuid.UUID `bun:"identifier,notnull,unique,type:uuid"`
	Name          string    `bun:"name,notnull"`
	Deleted       bool      `bun:"deleted,notnull,type:boolean,default:false"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:now()"`
	DeletedAt     time.Time `bun:"deleted_at,type:timestamp,default:null"`
}
