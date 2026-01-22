package client

import (
	"context"

	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
)

func (z *Zitadel) GetUserByUserName(ctx context.Context, username string) (*domain.User, error) {
	in := &management.GetUserByLoginNameGlobalRequest{LoginName: username}
	resp, err := z.api.ManagementService().GetUserByLoginNameGlobal(ctx, in)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		ZitadelID: resp.User.Id,
		UserName:  resp.User.UserName,
	}

	return u, nil
}

func (z *Zitadel) GetUserByID(ctx context.Context, userId string) (*domain.User, error) {
	in := &management.GetUserByIDRequest{
		Id: userId,
	}

	resp, err := z.api.ManagementService().GetUserByID(ctx, in)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		ZitadelID:    resp.User.Id,
		UserName:     resp.User.UserName,
		FirstName:    resp.User.GetHuman().Profile.FirstName,
		LastName:     resp.User.GetHuman().Profile.LastName,
		PrimaryEmail: resp.User.GetHuman().Email.Email,
		Role:         domain.OwnerRole,
	}

	return u, nil
}
