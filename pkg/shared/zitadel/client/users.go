package client

import (
	"context"

	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
)

func (z *Zitadel) GetUserByUserName(ctx context.Context, username string) (*domain.User, error) {
	inRequest := &user.ListUsersRequest{
		Queries: []*user.SearchQuery{
			{
				Query: &user.SearchQuery_UserNameQuery{
					UserNameQuery: &user.UserNameQuery{
						UserName: username,
					},
				},
			},
		},
	}

	resp, err := z.api.UserServiceV2().ListUsers(ctx, inRequest)
	if err != nil {
		return nil, err
	}

	userList := resp.GetResult()
	if len(userList) == 0 {
		return nil, nil
	}

	zitadelUser := userList[0]
	u := &domain.User{
		ZitadelID: zitadelUser.GetUserId(),
		UserName:  zitadelUser.GetUsername(),
	}

	return u, nil
}

func (z *Zitadel) GetUserByID(ctx context.Context, userId string) (*domain.User, error) {
	in := &user.GetUserByIDRequest{UserId: userId}

	resp, err := z.api.UserServiceV2().GetUserByID(ctx, in)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		ZitadelID:    resp.GetUser().GetUserId(),
		UserName:     resp.GetUser().GetUsername(),
		FirstName:    resp.GetUser().GetHuman().GetProfile().GetGivenName(),
		LastName:     resp.GetUser().GetHuman().GetProfile().GetFamilyName(),
		PrimaryEmail: resp.GetUser().GetHuman().GetEmail().GetEmail(),
		Role:         domain.OwnerRole,
	}

	return u, nil
}
