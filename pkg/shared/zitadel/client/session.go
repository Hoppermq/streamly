package client

import (
	"context"

	"github.com/hoppermq/middles"
	"github.com/hoppermq/streamly/pkg/domain/errors"
)

func (z *Zitadel) ValidateToken(ctx context.Context, token string) (*middles.Claims, error) {
	if z.verifier == nil {
		return nil, errors.ErrZitadelVerifierNotInitialized
	}

	ic, err := z.verifier.CheckAuthorization(ctx, token)
	if err != nil {
		return nil, err
	}

	claims := &middles.Claims{
		Subject:   ic.Subject,
		Email:     ic.Email,
		ExpiresAt: ic.ExpiresAt,
		IssuedAt:  ic.IssuedAt,
		Issuer:    ic.Issuer,
		Audience:  ic.Audience,
		OrgID:     ic.OrgID,
		Roles:     ic.Roles,
	}

	return claims, nil
}

func (z *Zitadel) RefreshToken(context.Context) (*middles.Claims, error) {
	return nil, errors.ErrZitadelInvalidRefreshToken
}
