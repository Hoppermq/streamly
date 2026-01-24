package ingestor

import (
	"context"
	"encoding/json"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/errors"
)

type EventIngestionUseCaseImpl struct {
	eventRepo domain.IngestionRepository

	logger *slog.Logger
	wg     sync.WaitGroup
}

type UseCaseOption func(*EventIngestionUseCaseImpl)

func WithEventRepository(eventRepo domain.IngestionRepository) UseCaseOption {
	return func(e *EventIngestionUseCaseImpl) {
		e.eventRepo = eventRepo
	}
}

func UseCaseWithLogger(logger *slog.Logger) UseCaseOption {
	return func(e *EventIngestionUseCaseImpl) {
		e.logger = logger
	}
}

func NewEventIngestionUseCase(opts ...UseCaseOption) domain.IngestionUseCase {
	useCase := &EventIngestionUseCaseImpl{
		wg: sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(useCase)
	}

	return useCase
}

func (uc *EventIngestionUseCaseImpl) IngestBatch(
	ctx context.Context,
	request *domain.BatchIngestionRequest,
) (*domain.BatchIngestionResponse, error) {
	uc.logger.Info("ingesting batch ingestion request", "request", request)
	var resp *domain.BatchIngestionResponse
	if err := uc.validateRequest(request); err != nil {
		return nil, errors.ErrEventCouldNotBeValidated
	}

	events, err := uc.transformToEvents(request)
	if err != nil {
		return nil, errors.ErrCouldNotTransformEvent
	}

	uc.wg.Add(1)
	go func(events []*domain.Event) {
		defer uc.wg.Done()
		if err := uc.eventRepo.BatchInsert(ctx, events); err != nil {
			uc.logger.Info("failed to ingest events", "error", err)
			resp = &domain.BatchIngestionResponse{
				Status:        "failed",
				IngestedCount: 0, // should be able to know
				Timestamp:     time.Now(),
				BatchID:       uuid.New().String(),
				FailedCount:   len(events), // should be able to know
			}
			return
		}
	}(events)
	uc.wg.Wait()

	if resp != nil {
		return resp, errors.ErrEventCouldNotInserted
	}

	resp = &domain.BatchIngestionResponse{
		Status:        "accepted",
		IngestedCount: len(events),
		Timestamp:     time.Now(),
		BatchID:       uuid.New().String(),
		FailedCount:   0,
	}

	uc.logger.Info("ingested batch ingestion response", "response", resp)
	return resp, nil
}

func (uc *EventIngestionUseCaseImpl) validateRequest(
	request *domain.BatchIngestionRequest,
) error {
	if request.TenantID == "" {
		return errors.ErrTenantIDRequired
	}
	if request.SourceID == "" {
		return errors.ErrSourceIDRequired
	}
	if request.Topic == "" {
		return errors.ErrSourceIDRequired
	}
	if len(request.Events) == 0 {
		return errors.ErrEventEmpty
	}
	if len(request.Events) > domain.EventBatchMaxSize {
		uc.logger.Info("events too big", "events", len(request.Events))
		return errors.ErrBatchSizeMaxSizeExceeded
	}

	for i := range request.Events {
		if err := uc.validateEventData(&request.Events[i], i); err != nil {
			return err
		}
	}

	return nil
}

func (uc *EventIngestionUseCaseImpl) validateEventData(
	event *domain.EventIngestionData,
	index int,
) error {
	if event.MessageID == "" {
		return errors.EventMessageMissing(index)
	}
	if event.EventType == "" {
		return errors.EventTypeMissing(index)
	}

	if len(event.Content) == 0 {
		return errors.EventContentEmpty(index)
	}

	var jsonContent interface{}
	if err := json.Unmarshal(event.Content, &jsonContent); err != nil {
		return errors.ErrEventCouldNotBeValidated
	}

	return nil
}

//nolint:gosec // false positive here ?
func (uc *EventIngestionUseCaseImpl) transformToEvents(
	request *domain.BatchIngestionRequest,
) ([]*domain.Event, error) {
	if len(request.Events) == 0 {
		return nil, errors.ErrEventEmpty
	}

	events := make([]*domain.Event, 0, len(request.Events))

	for i := range request.Events {
		eventData := &request.Events[i]
		if len(eventData.Content) > math.MaxUint32 {
			return nil, errors.ErrEventSize
		}

		event := &domain.Event{
			Timestamp:   time.Now(),
			TenantID:    request.TenantID,
			MessageID:   eventData.MessageID,
			SourceID:    request.SourceID,
			Topic:       request.Topic,
			ContentRaw:  string(eventData.Content),
			ContentJSON: eventData.Content,
			ContentSize: uint32(len(eventData.Content)),
			Headers:     eventData.Headers,
			FrameType:   eventData.FrameType,
			EventType:   eventData.EventType,
		}

		if event.Headers == nil {
			event.Headers = make(map[string]string)
		}

		events = append(events, event)
	}

	return events, nil
}
