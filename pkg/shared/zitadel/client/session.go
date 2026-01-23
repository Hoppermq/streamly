package client

import (
	"context"

	"github.com/hoppermq/middles"
)

func (z *Zitadel) ValidateToken(ctx context.Context, token string) (*middles.Claims, error) {
	return nil, nil
}

func (z *Zitadel) RefreshToken(context.Context) (*middles.Claims, error) {
	return nil, nil
}
