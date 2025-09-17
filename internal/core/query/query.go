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
	q.logger.Info("starting component", "name", q.serviceName, "service_id", q.serviceID)

	for _, handler := range q.handlers {
		q.wg.Add(1)

		go func(h domain.Service) {
			defer q.wg.Done()

			q.logger.Info("starting handler", "componnent", h.Name())
			if err := h.Run(ctx); err != nil {
				q.logger.WarnContext(ctx, "failed to run handler", "error", err)
				// should we stop if a service is not up ?
			}
		}(handler)
	}

	<-ctx.Done()
	q.logger.Info("stopping component", "name", q.serviceName, "service_id", q.serviceID)

	return nil

}

func (q *QueryService) Stop(ctx context.Context) error {
	q.logger.Info("stopping component", "name", q.serviceName, "service_id", q.serviceID)
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
	return q.serviceName
}

func (q *QueryService) SetServiceID(serviceID string) {
	q.serviceID = serviceID
}
