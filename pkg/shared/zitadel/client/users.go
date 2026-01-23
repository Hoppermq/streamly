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
		ZitadelID: resp.GetUser().GetId(),
		UserName:  resp.GetUser().GetUserName(),
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
		ZitadelID:    resp.GetUser().GetId(),
		UserName:     resp.GetUser().GetUserName(),
		FirstName:    resp.GetUser().GetHuman().GetProfile().GetFirstName(),
		LastName:     resp.GetUser().GetHuman().GetProfile().GetLastName(),
		PrimaryEmail: resp.GetUser().GetHuman().GetEmail().GetEmail(),
		Role:         domain.OwnerRole,
	}

	return u, nil
}
