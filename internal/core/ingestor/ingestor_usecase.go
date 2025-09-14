package ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/pkg/domain"
)

type EventIngestionUseCaseImpl struct {
	eventRepo domain.EventRepository

	logger *slog.Logger
}

type UseCaseOption func(*EventIngestionUseCaseImpl)

func WithEventRepository(eventRepo domain.EventRepository) UseCaseOption {
	return func(e *EventIngestionUseCaseImpl) {}
}

func UseCaseWithLogger(logger *slog.Logger) UseCaseOption {
	return func(e *EventIngestionUseCaseImpl) {
		e.logger = logger
	}
}

func NewEventIngestionUseCase(opts ...UseCaseOption) domain.EventIngestionUseCase {
	useCase := &EventIngestionUseCaseImpl{}

	for _, opt := range opts {
		opt(useCase)
	}

	return useCase
}

func (uc *EventIngestionUseCaseImpl) IngestBatch(ctx context.Context, request *domain.BatchIngestionRequest) (*domain.BatchIngestionResponse, error) {
	uc.logger.Info("ingesting batch ingestion request", "request", request)

	if err := uc.validateRequest(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	events, err := uc.transformToEvents(request)
	if err != nil {
		return nil, fmt.Errorf("transformation failed: %w", err)
	}

	if err := uc.eventRepo.BatchInsert(ctx, events); err != nil {
		uc.logger.Info("failed to ingest events", "error", err)
		return nil, fmt.Errorf("repository insert failed: %w", err)
	}

	response := &domain.BatchIngestionResponse{
		Status:        "accepted",
		IngestedCount: len(events),
		Timestamp:     time.Now(),
		BatchID:       uuid.New().String(),
		FailedCount:   0,
	}

	uc.logger.Info("ingested batch ingestion response", "response", response)
	return response, nil
}

func (uc *EventIngestionUseCaseImpl) validateRequest(request *domain.BatchIngestionRequest) error {
	if request.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if request.SourceID == "" {
		return fmt.Errorf("source_id is required")
	}
	if request.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if len(request.Events) == 0 {
		return fmt.Errorf("events cannot be empty")
	}
	if len(request.Events) > 5000 {
		uc.logger.Info("events too big", "events", len(request.Events))
		return fmt.Errorf("batch size cannot exceed 5000 events, got %d", len(request.Events))
	}

	for i, event := range request.Events {
		if err := uc.validateEventData(&event, i); err != nil {
			return err
		}
	}

	return nil
}

func (uc *EventIngestionUseCaseImpl) validateEventData(event *domain.EventIngestionData, index int) error {
	if event.MessageID == "" {
		return fmt.Errorf("events[%d]: message_id is required", index)
	}
	if event.EventType == "" {
		return fmt.Errorf("events[%d]: event_type is required", index)
	}
	if len(event.Content) == 0 {
		return fmt.Errorf("events[%d]: content cannot be empty", index)
	}

	var jsonContent interface{}
	if err := json.Unmarshal(event.Content, &jsonContent); err != nil {
		return fmt.Errorf("events[%d]: content must be valid JSON: %w", index, err)
	}

	return nil
}

func (uc *EventIngestionUseCaseImpl) transformToEvents(request *domain.BatchIngestionRequest) ([]*domain.Event, error) {
	events := make([]*domain.Event, 0, len(request.Events))

	for _, eventData := range request.Events {
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
