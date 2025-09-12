package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/http/routes"
	"github.com/hoppermq/streamly/pkg/domain"
)

// HTTPServer represent the HTTP server structure.
type HTTPServer struct {
	logger *slog.Logger

	health domain.HealthStatus

	engine *gin.Engine
	server *http.Server
}

type Options func(*HTTPServer)

// WithHTTPServer set the http server config.
func WithHTTPServer(conf *config.IngestionConfig) Options {
	return func(h *HTTPServer) {
		if h.engine == nil {
			panic(fmt.Errorf("engine should be set before server"))
		}
		h.server = &http.Server{
			Addr:         ":" + strconv.Itoa(conf.Ingestor.HTTP.Port),
			Handler:      h.engine,
			ReadTimeout:  conf.Ingestor.HTTP.ReadTimeout * time.Millisecond,
			WriteTimeout: conf.Ingestor.HTTP.WriteTimeout * time.Millisecond,
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
	routes.RegisterBaseRoutes(s.engine)

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Warn("http server stopped", "error", err)
		}
	}()

	return nil
}

// Shutdown stop gracefully the HTTPServer.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return nil
}

// IsHealthy return the health toolkit of the component.
func (s *HTTPServer) IsHealthy() bool {
	return true
}

func (s *HTTPServer) Name() string {
	return "http_handler"
}
