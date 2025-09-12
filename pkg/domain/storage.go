package domain

import "context"

type Storage interface {
	Save(ctx context.Context) error
	Get(ctx context.Context) error
}
