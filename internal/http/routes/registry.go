package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/pkg/domain"
)

// RouteRegistry manages route registration for different services.
type RouteRegistry struct {
	logger           *slog.Logger
	ingestionUseCase domain.IngestionUseCase
	// Future: queryUseCase domain.QueryUseCase
}

// RouteOption configures the route registry.
type RouteOption func(*RouteRegistry)

// WithLogger sets the logger for route registration.
func WithLogger(logger *slog.Logger) RouteOption {
	return func(r *RouteRegistry) {
		r.logger = logger
	}
}

// WithIngestionUseCase sets the ingestion use case for route registration.
func WithIngestionUseCase(useCase domain.IngestionUseCase) RouteOption {
	return func(r *RouteRegistry) {
		r.ingestionUseCase = useCase
	}
}

// NewRouteRegistry creates a new route registry with options.
func NewRouteRegistry(opts ...RouteOption) *RouteRegistry {
	r := &RouteRegistry{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// RegisterAllRoutes registers all configured routes.
func (r *RouteRegistry) RegisterAllRoutes(engine *gin.Engine) {
	RegisterBaseRoutes(engine)

	if r.ingestionUseCase != nil && r.logger != nil {
		RegisterIngestionRoutes(engine, r.logger, r.ingestionUseCase)
	}
}

// CreateRouteRegistrar creates a route registrar function with the given options.
func CreateRouteRegistrar(opts ...RouteOption) func(*gin.Engine) {
	return func(engine *gin.Engine) {
		registry := NewRouteRegistry(opts...)
		registry.RegisterAllRoutes(engine)
	}
}
