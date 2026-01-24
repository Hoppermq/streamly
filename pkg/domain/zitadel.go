package domain

import (
	"context"

	"github.com/google/uuid"
)

type Client interface {
	GetUserByUserName(ctx context.Context, userName string) (*User, error)
	GetUserByID(ctx context.Context, userId string) (*User, error)

	GetOrganizationByID(ctx context.Context, organizationId uuid.UUID) (*Organization, error)
}

//nolint:tagliatelle
type ZitadelOrganization struct {
	OrganizationID string `binding:"required" json:"orgId"`
}

//nolint:tagliatelle
type ZitadelProfile struct {
	FirstName string `binding:"required" json:"givenName"`
	LastName  string `binding:"required" json:"familyName"`
}

//nolint:tagliatelle
type ZitadelEmail struct {
	Email      string `binding:"required,email" json:"email"`
	IsVerified bool   `binding:"required"       json:"isVerified"`
}

//nolint:tagliatelle
type ZitadelEventUserCreatedRequest struct {
	Email        ZitadelEmail        `binding:"required" json:"email"`
	Organization ZitadelOrganization `binding:"required" json:"organization"`
	Profile      ZitadelProfile      `binding:"required" json:"profile"`
	UserName     string              `binding:"required" json:"userName"`
}

//nolint:tagliatelle
type ZitadelEventUserCreated struct {
	InstanceID     string                         `binding:"required" json:"instanceId"`
	OrganizationID string                         `binding:"required" json:"orgID"`
	Request        ZitadelEventUserCreatedRequest `binding:"required" json:"request"`
	UserID         string                         `binding:"required" json:"userID"`
}
