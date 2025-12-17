package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Organization struct {
	bun.BaseModel `bun:"table:tenants,alias:org"`
	ID            string    `bun:"id,pk,type:uuid,scanonly"`
	Identifier    string    `bun:"identifier,notnull,unique,type:uuid"`
	Name          string    `bun:"name,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:now()"`
}
