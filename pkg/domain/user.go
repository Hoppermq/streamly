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
	Role         PlatformRole

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	WithTx(tx TxContext) UserRepository

	FindOneByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindOneByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, limit, offset int) ([]User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exist(ctx context.Context, identifier uuid.UUID) (bool, error)
	GetUserIDFromZitadelID(ctx context.Context, zitadelID string) (uuid.UUID, error)
}

type CreateUser struct {
	UserName     string `binding:"required" form:"user_name"     json:"user_name"`
	FirstName    string `binding:"required" form:"first_name"    json:"first_name"`
	LastName     string `binding:"required" form:"last_name"     json:"last_name"`
	PrimaryEmail string `binding:"required" form:"primary_email" json:"primary_email"`
	Role         string `form:"role"        json:"role"`
	ZitadelID    string `form:"zitadel_id"  json:"zitadel_id"`
}

type UpdateUser struct {
	UserName     string
	FirstName    string
	LastName     string
	PrimaryEmail string
	Role         string
}
