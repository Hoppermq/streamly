package http

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/pkg/domain/errors"
)

// Server represent the HTTP server structure.
type Server struct {
	logger *slog.Logger

	engine *gin.Engine
	server *http.Server
}

type Options func(*Server)

// WithHTTPServer set the http server config for ingestion service.
func WithHTTPServer(conf *config.IngestionConfig) Options {
	return func(h *Server) {
		if h.engine == nil {
			panic(errors.ErrEngineErrorOrder)
		}
		h.server = &http.Server{
			Addr:         ":" + strconv.Itoa(conf.Ingestor.HTTP.Port),
			Handler:      h.engine,
			ReadTimeout:  conf.Ingestor.HTTP.ReadTimeout,
			WriteTimeout: conf.Ingestor.HTTP.WriteTimeout,
		}
	}
}

// WithQueryHTTPServer set the http server config for query service.
func WithQueryHTTPServer(conf *config.QueryConfig) Options {
	return func(h *Server) {
		if h.engine == nil {
			panic(errors.ErrEngineErrorOrder)
		}
		h.server = &http.Server{
			Addr:         ":" + strconv.Itoa(conf.Query.HTTP.Port),
			Handler:      h.engine,
			ReadTimeout:  conf.Query.HTTP.ReadTimeout,
			WriteTimeout: conf.Query.HTTP.WriteTimeout,
		}
	}
}

func WithAuthHTTPServer(conf *config.AuthConfig) Options {
	return func(h *Server) {
		if h.engine == nil {
			panic(errors.ErrEngineErrorOrder)
		}
		h.server = &http.Server{
			Addr:         ":" + strconv.Itoa(conf.Auth.HTTP.Port),
			Handler:      h.engine,
			ReadTimeout:  conf.Auth.HTTP.ReadTimeout,
			WriteTimeout: conf.Auth.HTTP.WriteTimeout,
		}
	}
}

func WithPlatformHTTPServer(conf *config.PlatformConfig) Options {
	return func(h *Server) {
		if h.engine == nil {
			panic(errors.ErrEngineErrorOrder)
		}
		h.server = &http.Server{
			Addr:         ":" + strconv.Itoa(conf.Platform.HTTP.Port),
			Handler:      h.engine,
			ReadTimeout:  conf.Platform.HTTP.ReadTimeout,
			WriteTimeout: conf.Platform.HTTP.WriteTimeout,
		}
	}
}

// WithEngine set the engine.
func WithEngine(engine *gin.Engine) Options {
	return func(h *Server) {
		h.engine = engine
	}
}

// WithLogger set the logger.
func WithLogger(logger *slog.Logger) Options {
	return func(h *Server) {
		h.logger = logger
	}
}

// WithRoutes register all routes to the engine.
func WithRoutes(routerRegister func(engine *gin.Engine)) Options {
	return func(h *Server) {
		routerRegister(h.engine)
	}
}

// NewHTTPServer create the http server.
func NewHTTPServer(opts ...Options) *Server {
	h := &Server{}
	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Run will run the http server component.
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("starting http server", "addr", s.server.Addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Warn("http server stopped", "error", err)
		}
	}()

	return nil
}

// Shutdown stop gracefully the HTTPServer.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Warn("http server shutdown failed", "error", err)
		return err
	}

	return nil
}

// IsHealthy return the health toolkit of the component.
func (s *Server) IsHealthy() bool {
	return true
}

func (s *Server) Name() string {
	return "http_handler"
}
