package platform

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/hoppermq/streamly/pkg/domain"
)

type Streamly struct {
	logger   *slog.Logger
	handlers []domain.Service

	serviceID   string
	serviceName string

	wg     *sync.WaitGroup
	mu     sync.RWMutex
	cancel context.CancelFunc
}

type Option func(*Streamly) error

func WithLogger(logger *slog.Logger) Option {
	return func(s *Streamly) error {
		s.logger = logger
		return nil
	}
}

func WithHandler(handlers ...domain.Service) Option {
	return func(s *Streamly) error {
		s.handlers = handlers
		return nil
	}
}

func NewStreamlyService(opts ...Option) *Streamly {
	streamly := &Streamly{
		wg: new(sync.WaitGroup),
		mu: sync.RWMutex{},
	}

	for _, opt := range opts {
		err := opt(streamly)
		if err != nil {
			return nil
		}
	}

	return streamly
}

func (s *Streamly) Run(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	s.logger.InfoContext(ctx, "starting component", "name", s.serviceName, "service_id", s.serviceID)

	for _, handler := range s.handlers {
		s.wg.Add(1)

		go func(h domain.Service) {
			defer s.wg.Done()
			s.logger.InfoContext(ctx, "starting handler", "component", h.Name())
			if err := h.Run(ctx); err != nil {
				s.logger.WarnContext(ctx, "failed to start handler", "component", h.Name())
			}
		}(handler)
	}

	<-ctx.Done()
	s.logger.Info("shutting down component", "service_name", s.serviceName, "service_id", s.serviceID)
	return nil
}

func (s *Streamly) Stop(ctx context.Context) error {
	s.logger.InfoContext(ctx, "stopping component", "name", s.serviceName, "service_id", s.serviceID)
	if s.cancel != nil {
		s.cancel()
	}

	var errs []error
	for _, handler := range s.handlers {
		if err := handler.Shutdown(ctx); err != nil {
			s.logger.WarnContext(ctx, "failed to stop handler", "component", handler.Name())
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.ErrUnsupported
	}

	s.wg.Wait()
	return nil
}

func (s *Streamly) Name() string {
	return s.serviceName
}

func (s *Streamly) SetServiceID(serviceID string) { // Should be set on init !
	s.mu.Lock()
	defer s.mu.Unlock()

	s.serviceID = serviceID
}
