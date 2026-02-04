package bootstrap

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/errors"
)

func (o *Orchestrator) isFirstInstance(ctx context.Context) (bool, error) {
	rootUserMail := os.Getenv("ROOT_USER_EMAIL")

	if rootUserMail == "" {
		o.logger.WarnContext(
			ctx,
			"ROOT_USER_EMAIL environment variable not set",
			"error",
			errors.ErrRootEmailNotSet,
		)
		rootUserMail = "root@streamly.auth.localhost"
		o.logger.Info("ROOT_USER_EMAIL switched to default")
	}

	// will use env variable here
	owner, err := o.userUC.FindOneByPrimaryEmail(ctx, rootUserMail)
	if err != nil && !isNotFoundError(err) {
		return false, errors.RootUserQueryFailed(err, rootUserMail)
	}

	if owner != nil {
		return false, nil
	}

	localOrg, err := o.orgUC.FindAll(ctx, 1, 0)
	if err != nil && !isNotFoundError(err) {
		return false, fmt.Errorf("failed to query organizations: %w", err)
	}

	if len(localOrg) > 0 {
		return false, nil
	}

	return true, nil
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "sql: no rows in result set" || strings.Contains(err.Error(), "no rows")
}

func (o *Orchestrator) setupDefaultOrg(ctx context.Context) error {
	rootUserID := os.Getenv("ROOT_USER_ID")
	if rootUserID == "" {
		return errors.ErrRootUserIDNotSet
	}

	u, err := o.zitadel.GetUserByID(ctx, rootUserID)
	if err != nil {
		return errors.RootUserQueryFailed(err, rootUserID)
	}

	createUserCommand := domain.CreateUser{
		UserName:     u.UserName,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		ZitadelID:    u.ZitadelID,
		PrimaryEmail: u.PrimaryEmail,

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
