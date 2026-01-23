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

type ZitadelOrganization struct {
	OrganizationID string `binding:"required" json:"orgId"`
}

type ZitadelProfile struct {
	FirstName string `binding:"required" json:"givenName"`
	LastName  string `binding:"required" json:"familyName"`
}

type ZitadelEmail struct {
	Email      string `binding:"required,email" json:"email"`
	IsVerified bool   `binding:"required"       json:"isVerified"`
}

type ZitadelEventUserCreatedRequest struct {
	Email        ZitadelEmail        `binding:"required" json:"email"`
	Organization ZitadelOrganization `binding:"required" json:"organization"`
	Profile      ZitadelProfile      `binding:"required" json:"profile"`
	UserName     string              `binding:"required" json:"userName"`
}
type ZitadelEventUserCreated struct {
	InstanceID     string                         `binding:"required" json:"instanceId"`
	OrganizationID string                         `binding:"required" json:"orgID"`
	Request        ZitadelEventUserCreatedRequest `binding:"required" json:"request"`
	UserID         string                         `binding:"required" json:"userID"`
}
