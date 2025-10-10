package query

import (
	"context"
	"log/slog"
	"sync"

	"github.com/hoppermq/streamly/pkg/domain"
)

type QueryService struct {
	logger *slog.Logger

	handlers []domain.Service

	serviceID   string
	serviceName string

	mu     sync.Mutex
	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

type Option func(*QueryService)

func NewQueryService(opts ...Option) *QueryService {
	q := &QueryService{
		mu: sync.Mutex{},
		wg: &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt(q)
	}

	return q
}

func WithLogger(logger *slog.Logger) Option {
	return func(q *QueryService) {
		q.logger = logger
	}
}

func WithHandlers(handlers ...domain.Service) Option {
	return func(q *QueryService) {
		q.handlers = handlers
	}
}

func (q *QueryService) Run(ctx context.Context) error {
	ctx, q.cancel = context.WithCancel(ctx)

	q.mu.Lock()
	serviceName := q.serviceName
	serviceID := q.serviceID
	q.mu.Unlock()

	q.logger.Info("starting component", "name", serviceName, "service_id", serviceID)

	for _, handler := range q.handlers {
		q.wg.Add(1)

		go func(h domain.Service) {
			defer q.wg.Done()

			q.logger.Info("starting handler", "componnent", h.Name())
			if err := h.Run(ctx); err != nil {
				q.logger.WarnContext(ctx, "failed to run handler", "error", err)
			}
		}(handler)
	}

	<-ctx.Done()

	q.mu.Lock()
	serviceName = q.serviceName
	serviceID = q.serviceID
	q.mu.Unlock()

	q.logger.Info("stopping component", "name", serviceName, "service_id", serviceID)

	return nil

}

func (q *QueryService) Stop(ctx context.Context) error {
	q.mu.Lock()
	serviceName := q.serviceName
	serviceID := q.serviceID
	q.mu.Unlock()

	q.logger.Info("stopping component", "name", serviceName, "service_id", serviceID)
	if q.cancel != nil {
		q.cancel()
	}

	for _, handler := range q.handlers {
		if err := handler.Shutdown(ctx); err != nil {
			q.logger.WarnContext(ctx, "failed to shutdown handler", "error", err)
			return err
		}
		q.logger.Info("shutdown handler", "componnent", handler.Name())
	}

	q.wg.Wait()
	return nil
}

func (q *QueryService) Name() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.serviceName
}

func (q *QueryService) SetServiceID(serviceID string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.serviceID = serviceID
}
