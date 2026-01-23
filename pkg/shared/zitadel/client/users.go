package client

import (
	"context"

	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
)

func (z *Zitadel) GetUserByUserName(ctx context.Context, username string) (*domain.User, error) {
	in := &user.ListUsersRequest{
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

	resp, err := z.api.UserServiceV2().ListUsers(ctx, in)
	if err != nil {
		return nil, err
	}

	userList := resp.GetResult()
	if len(userList) == 0 {
		return nil, nil
	}

	zitadelUser := userList[0]
	u := &domain.User{
		ZitadelID: zitadelUser.UserId,
		UserName:  zitadelUser.Username,
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
		ZitadelID:    resp.GetUser().UserId,
		UserName:     resp.GetUser().Username,
		FirstName:    resp.GetUser().GetHuman().GetProfile().GivenName,
		LastName:     resp.GetUser().GetHuman().GetProfile().FamilyName,
		PrimaryEmail: resp.GetUser().GetHuman().GetEmail().Email,
		Role:         domain.OwnerRole,
	}

	return u, nil
}
