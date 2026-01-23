package client

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/hoppermq/streamly/pkg/domain/errors"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

type Zitadel struct {
	logger *slog.Logger
	api    *client.Client
	pat    string
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

	if zita.logger != nil {
		zita.logger.InfoContext(ctx, "Zitadel client initialized successfully")
	}

	return zita, nil
}
