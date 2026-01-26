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

	if z.cache != nil {
		if cached, ok := z.cache.Get(token); ok {
			return cached, nil
		}
	}

	resp, err := z.authVerifier.CheckAuthorization(ctx, token)

	if err != nil {
		return nil, err
	}

	claims := &middles.Claims{
		Subject:   resp.Subject,
		Email:     resp.Email,
		ExpiresAt: resp.Expiration.AsTime().Unix(),
		IssuedAt:  resp.IssuedAt.AsTime().Unix(),
		Issuer:    resp.Issuer,
		Audience:  resp.Audience,
		OrgID:     resp.OrganizationID(),
	}

	if z.cache != nil {
		z.cache.Set(token, claims, z.cacheTTL)
	}

	return claims, nil
}

func (z *Zitadel) RefreshToken(context.Context) (*middles.Claims, error) {
	return nil, errors.ErrZitadelInvalidRefreshToken
}
