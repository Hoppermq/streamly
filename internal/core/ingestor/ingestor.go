// Package ingestor represent the ingestion service package.
package ingestor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/pkg/domain"
	serviceloader "github.com/zixyos/goloader/service"
)

// Ingestor represent the ingestion service component.
type Ingestor struct {
	logger *slog.Logger

	handlers []serviceloader.Service // this will be the domain type for handlers

	serviceID   string
	serviceName string

	wg     *sync.WaitGroup
	mu     sync.Mutex
	cancel context.CancelFunc
}

// Option is the function callback to compose the Ingestor.
type Option func(*Ingestor)

func WithLogger(logger *slog.Logger) Option {
	return func(ingestor *Ingestor) {
		ingestor.logger = logger
	}
}

// WithHandlers load the handlers to the service.
func WithHandlers(handlers ...serviceloader.Service) Option {
	return func(ingestor *Ingestor) {
		ingestor.handlers = handlers
	}
}

// WithServiceName load service name.
func WithServiceName(name string) Option {
	return func(ingestor *Ingestor) {
		ingestor.serviceName = name
	}
}

// WithServiceID load the serviceID.
func WithServiceID(id string) Option {
	return func(ingestor *Ingestor) {
		ingestor.serviceID = id
	}
}

// Run will start the component.
func (i Ingestor) Run(ctx context.Context) error {
	ctx, i.cancel = context.WithCancel(ctx)

	i.logger.Info(
		"starting component",
		"name",
		i.serviceName,
		"id",
		i.serviceID,
	)

	for _, handler := range i.handlers {
		i.wg.Add(1)

		go func(h serviceloader.Service) {
			defer i.wg.Done()
			i.logger.Info("starting handler", "component", handler.Name())

			if err := h.Run(ctx); err != nil {
				i.logger.Error("handler failed", "name", h.Name(), "error", err)
			}
		}(handler)
	}

	<-ctx.Done()
	i.logger.Info("service is shutting down")

	return nil
}

// Stop shutdown gracefully all
func (i *Ingestor) Stop(ctx context.Context) error {
	i.logger.Info("stopping component", "name", i.serviceName)

	if i.cancel != nil {
		i.cancel()
	}

	for _, handler := range i.handlers {
		if err := handler.Stop(ctx); err != nil {
			i.logger.Warn("failed to shutdown handler", "name", handler.Name(), "error", err)
			return err
		}
		i.logger.Info("shutdown handler successfully", "name", handler.Name())
	}

	i.wg.Wait()
	return nil
}

// SetServiceID load the service id.
func (i *Ingestor) SetServiceID(serviceID string) {
	defer i.mu.Unlock()

	i.mu.Lock()
	i.serviceID = serviceID
}

// Name return the service name.
func (i *Ingestor) Name() string {
	return i.serviceName
}

func WithConfig(cfg *config.IngestionConfig) Option {
	return func(ingestor *Ingestor) {
		ingestor.serviceName = cfg.Ingestor.Service.Name
	}
}

func WithTransport(transport domain.Transport) Option {
	return func(ingestor *Ingestor) {

	}
}

func NewIngestor(opts ...Option) (*Ingestor, error) {
	ingestor := &Ingestor{
		wg: new(sync.WaitGroup),
		mu: sync.Mutex{},
	}

	for _, opt := range opts {
		opt(ingestor)
	}

	return ingestor, nil
}
