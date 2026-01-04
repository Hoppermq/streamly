package domain

import (
	"context"

	"github.com/google/uuid"
)

type Client interface {
	GetUserByUserName(ctx context.Context, userName string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, userId string) (*User, error)

	GetOrganizationByID(ctx context.Context, organizationId uuid.UUID) (*Organization, error)
}
type ZitadelEventUserCreatedRequest struct {
	Email struct {
		Email      string `json:"email" binding:"required,email"`
		IsVerified bool   `json:"isVerified" binding:"required"`
	} `json:"email" binding:"required"`
	Organization struct {
		OrganizationID string `json:"orgId" binding:"required"`
	} `json:"organization" binding:"required"`
	Profile struct {
		FirstName string `binding:"required" json:"givenName"`
		LastName  string `binding:"required" json:"familyName"`
	} `json:"profile" binding:"required"`
	UserName string `json:"userName" binding:"required"`
}
type ZitadelEventUserCreated struct {
	InstanceID     string                         `json:"instanceId" binding:"required"`
	OrganizationID string                         `json:"orgID" binding:"required"`
	Request        ZitadelEventUserCreatedRequest `json:"request" binding:"required"`
	UserID         string                         `json:"userID" binding:"required"`
}
