package domain

import "context"

// Service represent the domain type of service component
type Service interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Name() string
	HealthStatus
}

// HealthStatus represent the domain type for service with health
type HealthStatus interface {
	IsHealthy() bool
}
