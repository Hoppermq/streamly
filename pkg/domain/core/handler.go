package domain

import "context"

type HealthCallback func(ctx context.Context) (bool, error)
