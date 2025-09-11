package domain

import "context"

type Service interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Health() HealthStatus
}

type HealthStatus interface {
	IsHealthy() bool
}
