package models

import "github.com/uptrace/bun"

type UserAuth struct {
	bun.BaseModel `bun:"table:projections.users14_humans"`
	ID            string `bun:"user_id,type:string,notnull"`
	InstanceID    string `bun:"instance_id,type:string,notnull"`
	FirstName     string `bun:"first_name,type:string,notnull"`
	LastName      string `bun:"last_name,type:string,notnull"`
	Email         string `bun:"email,type:string,notnull"`
	DisplayName   string `bun:"display_name,type:string"`
}

type UserProjection struct {
	bun.BaseModel `bun:"table:projections.users14"`
	ID            string `bun:"id,type:string,notnull"`
	InstanceID    string `bun:"instance_id,type:string,notnull"`
	UserName      string `bun:"username,type:string,notnull"`
}
