package bootstrap

import (
	"context"
	"fmt"
	"os"

	"github.com/hoppermq/streamly/pkg/domain"
)

func (o *Orchestrator) setupDefaultOrg(ctx context.Context) error {
	rootUserID := os.Getenv("ROOT_USER_ID")
	if rootUserID == "" {
		return fmt.Errorf("ROOT_USER_ID environment variable is not set")
	}

	u, err := o.zitadel.GetUserByID(ctx, rootUserID)
	if err != nil {
		return fmt.Errorf("failed to get root user by ID %s: %w", rootUserID, err)
	}

	createUserCommand := domain.CreateUser{
		UserName:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ZitadelID: u.ZitadelID,

		Role: string(u.Role),
	}

	if err := o.userUC.Create(ctx, &createUserCommand); err != nil {
		return err
	}

	createOrganizationCommand := domain.CreateOrganization{
		Name: "local",
	}

	if err := o.orgUC.Create(ctx, createOrganizationCommand, u.ZitadelID); err != nil {
		return err
	}

	return nil
}
