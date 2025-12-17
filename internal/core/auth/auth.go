package auth

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/hoppermq/streamly/pkg/domain"
)

type Auth struct {
	logger *slog.Logger

	handlers []domain.Service

	serviceID   string
	serviceName string

	wg     *sync.WaitGroup
	mu     sync.RWMutex
	cancel context.CancelFunc
}

type Option func(*Auth)

func WithLogger(logger *slog.Logger) Option {
	return func(a *Auth) {
		a.logger = logger
	}
}

func WithHandler(handlers ...domain.Service) Option {
	return func(a *Auth) {
		a.handlers = handlers
	}
}

func NewAuthService(opts ...Option) *Auth {
	auth := &Auth{
		wg: new(sync.WaitGroup),
		mu: sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(auth)
	}

	return auth
}

func (a *Auth) Run(ctx context.Context) error {
	ctx, a.cancel = context.WithCancel(ctx)

	a.logger.Info("starting component", "name", a.serviceName, "service_id", a.serviceID)

	for _, handler := range a.handlers {
		a.wg.Add(1)

		go func(h domain.Service) {
			defer a.wg.Done()
			a.logger.Info("starting handler", "component", h.Name())
			if err := h.Run(ctx); err != nil {
				a.logger.Warn("failed to start handler", "component", h.Name(), "err", err)
			}
		}(handler)
	}
	<-ctx.Done()
	a.logger.Info("service is shutting down")
	return nil
}

func (a *Auth) Stop(ctx context.Context) error {
	a.logger.Info("shutting down component", "name", a.serviceName, "service_id", a.serviceID)
	if a.cancel != nil {
		a.cancel()
	}

	var errs []error
	for _, handler := range a.handlers {
		if err := handler.Shutdown(ctx); err != nil {
			a.logger.Warn("failed to shutdown handler", "component", handler.Name(), "err", err)
			errs = append(errs, err)
		}
		a.logger.Info("shutdown handler", "component", handler.Name())
	}

	if len(errs) > 0 {
		return errors.ErrUnsupported
	}

	a.wg.Wait()
	return nil
}

func (a *Auth) Name() string {
	return a.serviceName
}

func (a *Auth) SetServiceID(serviceID string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.serviceID = serviceID
}
