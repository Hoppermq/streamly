package http

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/errors"
)

// HTTPServer represent the HTTP server structure.
type HTTPServer struct {
	logger *slog.Logger

	health domain.HealthStatus

	engine *gin.Engine
	server *http.Server
}

type Options func(*HTTPServer)

// WithHTTPServer set the http server config for ingestion service.
func WithHTTPServer(conf *config.IngestionConfig) Options {
	return func(h *HTTPServer) {
		if h.engine == nil {
			panic(errors.ErrEngineErrorOrder)
		}
		h.server = &http.Server{
			Addr:         ":" + strconv.Itoa(conf.Ingestor.HTTP.Port),
			Handler:      h.engine,
			ReadTimeout:  conf.Ingestor.HTTP.ReadTimeout * time.Millisecond,
			WriteTimeout: conf.Ingestor.HTTP.WriteTimeout * time.Millisecond,
		}
	}
}

// WithQueryHTTPServer set the http server config for query service.
func WithQueryHTTPServer(conf *config.QueryConfig) Options {
	return func(h *HTTPServer) {
		if h.engine == nil {
			panic(errors.ErrEngineErrorOrder)
		}
		h.server = &http.Server{
			Addr:         ":" + strconv.Itoa(conf.Query.HTTP.Port),
			Handler:      h.engine,
			ReadTimeout:  conf.Query.HTTP.ReadTimeout * time.Millisecond,
			WriteTimeout: conf.Query.HTTP.WriteTimeout * time.Millisecond,
		}
	}
}

func WithAuthHTTPServer(conf *config.AuthConfig) Options {
	return func(h *HTTPServer) {
		if h.engine == nil {
			panic(errors.ErrEngineErrorOrder)
		}
		h.server = &http.Server{
			Addr:         ":" + strconv.Itoa(conf.Auth.HTTP.Port),
			Handler:      h.engine,
			ReadTimeout:  conf.Auth.HTTP.ReadTimeout * time.Millisecond,
			WriteTimeout: conf.Auth.HTTP.WriteTimeout * time.Millisecond,
		}
	}
}

// WithEngine set the engine.
func WithEngine(engine *gin.Engine) Options {
	return func(h *HTTPServer) {
		h.engine = engine
	}
}

// WithLogger set the logger.
func WithLogger(logger *slog.Logger) Options {
	return func(h *HTTPServer) {
		h.logger = logger
	}
}

// WithRoutes register all routes to the engine.
func WithRoutes(routerRegister func(engine *gin.Engine)) Options {
	return func(h *HTTPServer) {
		routerRegister(h.engine)
	}
}

// NewHTTPServer create the http server.
func NewHTTPServer(opts ...Options) *HTTPServer {
	h := &HTTPServer{}
	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Run will run the http server component.
func (s *HTTPServer) Run(ctx context.Context) error {
	s.logger.Info("starting http server", "addr", s.server.Addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Warn("http server stopped", "error", err)
		}
	}()

	return nil
}

// Shutdown stop gracefully the HTTPServer.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Warn("http server shutdown failed", "error", err)
		return err
	}

	return nil
}

// IsHealthy return the health toolkit of the component.
func (s *HTTPServer) IsHealthy() bool {
	return true
}

func (s *HTTPServer) Name() string {
	return "http_handler"
}
