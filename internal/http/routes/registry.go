package routes

import (
	"github.com/gin-gonic/gin"
)

// RouteRegistrar defines a function that registers routes on an engine.
type RouteRegistrar func(*gin.Engine)

// RouteRegistry manages multiple route registrars.
type RouteRegistry struct {
	registrars []RouteRegistrar
}

// NewRouteRegistry creates a new route registry.
func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		registrars: make([]RouteRegistrar, 0),
	}
}

// AddRegistrar adds a route registrar to the registry.
func (r *RouteRegistry) AddRegistrar(registrar RouteRegistrar) {
	r.registrars = append(r.registrars, registrar)
}

// RegisterAll registers all routes with the provided engine.
func (r *RouteRegistry) RegisterAll(engine *gin.Engine) {
	RegisterBaseRoutes(engine)

	for _, registrar := range r.registrars {
		registrar(engine)
	}
}

// CreateRouteRegistrar creates a route registrar function from multiple registrars.
func CreateRouteRegistrar(registrars ...RouteRegistrar) func(*gin.Engine) {
	return func(engine *gin.Engine) {
		registry := NewRouteRegistry()
		for _, registrar := range registrars {
			registry.AddRegistrar(registrar)
		}
		registry.RegisterAll(engine)
	}
}
