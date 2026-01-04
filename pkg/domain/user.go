package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Identifier uuid.UUID
	ZitadelID  string

	UserName     string
	FirstName    string
	LastName     string
	PrimaryEmail string
	Role         string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	FindOneByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindAll(ctx context.Context, limit, offset int) ([]User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateUser struct {
	UserName     string `form:"user_name" json:"user_name" binding:"required"`
	FirstName    string `form:"first_name" json:"first_name" binding:"required"`
	LastName     string `form:"last_name" json:"last_name" binding:"required"`
	PrimaryEmail string `form:"primary_email" json:"primary_email" binding:"required"`
	Role         string `form:"role" json:"role"`
	ZitadelID    string `form:"zitadel_id" json:"zitadel_id"`
}

type UpdateUser struct {
	UserName     string
	FirstName    string
	LastName     string
	PrimaryEmail string
	Role         string
}
