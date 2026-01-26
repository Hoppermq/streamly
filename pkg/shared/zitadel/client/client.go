package client

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"

	"github.com/hoppermq/middles"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/errors"
)

const defaultMaxTime = 5 * time.Minute

type tokenVerifier interface {
	CheckAuthorization(ctx context.Context, token string) (*oauth.IntrospectionContext, error)
}
type Zitadel struct {
	logger *slog.Logger
	api    *client.Client
	pat    string

	verifier     *oauth.IntrospectionVerificationWithCache[*middles.Claims]
	authVerifier tokenVerifier
	cache        domain.Cache[*middles.Claims]
	cacheTTL     time.Duration
	keyfilePath  string
	issuer       string
}

type Options func(*Zitadel) error

func WithLogger(logger *slog.Logger) Options {
	return func(z *Zitadel) error {
		z.logger = logger
		return nil
	}
}

// WithPAT sets the Personal Access Token directly (for prod: loaded from env/SSM).
func WithPAT(pat string) Options {
	return func(z *Zitadel) error {
		z.pat = strings.TrimSpace(pat)
		return nil
	}
}

// WithPATFromFile loads PAT from file path (for v0: local dev).
//
//nolint:gosec // TODO : use FS.
func WithPATFromFile(path string) Options {
	return func(z *Zitadel) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return errors.FailedToReadFile(path)
		}
		z.pat = strings.TrimSpace(string(data))
		return nil
	}
}

func WithPort(port uint16) zitadel.Option {
	return zitadel.WithPort(port)
}

func WithInsecure(port string) zitadel.Option {
	return zitadel.WithInsecure(port)
}

func NewZitadel(domain string, opts ...zitadel.Option) *zitadel.Zitadel {
	return zitadel.New(domain, opts...)
}

func WithTokenCache[T comparable](cache domain.Cache[*middles.Claims], ttl time.Duration) Options {
	return func(z *Zitadel) error {
		z.cache = cache
		z.cacheTTL = ttl
		return nil
	}
}

// WithServiceAccountKeyFile sets the keyfile path for ResourceServer authentication.
func WithServiceAccountKeyFile(path string) Options {
	return func(z *Zitadel) error {
		z.keyfilePath = strings.TrimSpace(path)
		return nil
	}
}

// WithIssuer sets the Zitadel issuer URL (e.g., "https://your-domain.zitadel.cloud").
func WithIssuer(issuer string) Options {
	return func(z *Zitadel) error {
		z.issuer = strings.TrimSpace(issuer)
		return nil
	}
}

// NewZitadelClient create a new instance of the zitadel api client.
func NewZitadelClient(ctx context.Context, z *zitadel.Zitadel, opts ...Options) (*Zitadel, error) {
	zita := &Zitadel{}

	for _, opt := range opts {
		if err := opt(zita); err != nil {
			return nil, err
		}
	}

	if zita.pat == "" {
		return nil, errors.ErrZitadelPATRequired
	}

	authOptions := client.PAT(zita.pat)
	c, err := client.New(ctx, z, client.WithAuth(authOptions))

	if err != nil {
		return nil, errors.ErrZitadelClientCreation
	}

	zita.api = c

	resourceServer, err := rs.NewResourceServerFromKeyFile(ctx, zita.issuer, zita.keyfilePath)
	if err != nil {
		return nil, errors.ZitadelResourceServerCreationFailed(err)
	}

	ttl := zita.cacheTTL
	if ttl == 0 {
		ttl = defaultMaxTime
	}

	zita.verifier = oauth.NewIntrospectionVerificationWithCache(resourceServer, zita.cache, ttl)

	// TODO : bootstrap need to create the platform sa.
	v2 := oauth.DefaultJWTAuthorization("356531635715399939")
	verifier, err := v2(ctx, z)
	if err != nil {
		return nil, err
	}

	zita.authVerifier = verifier

	if zita.logger != nil {
		zita.logger.InfoContext(ctx, "Zitadel client initialized successfully")
	}

	return zita, nil
}
