package ingestor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hoppermq/streamly/pkg/domain"
)

type EventRepository struct {
	driver domain.Driver
}

type RepositoryOption func(*EventRepository)

func WithDriver(driver domain.Driver) RepositoryOption {
	return func(r *EventRepository) {
		r.driver = driver
	}
}

func (e EventRepository) BatchInsert(ctx context.Context, events []*domain.Event) error {
	return nil
}

func NewEventRepository(opts ...RepositoryOption) *EventRepository {
	r := &EventRepository{}
	for _, opt := range opts {
		opt(r)
	}

	return r
}

type MockEventRepository struct {
	events      []domain.Event
	failureRate float64
}

func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{
		events:      make([]domain.Event, 0),
		failureRate: 0.0,
	}
}

func (r *MockEventRepository) BatchInsert(ctx context.Context, events []*domain.Event) error {
	log.Printf("MockEventRepository: Simulating batch insert of %d events", len(events))

	for i, event := range events {
		if err := r.validateEvent(event); err != nil {
			return fmt.Errorf("validation failed for event %d: %w", i, err)
		}
	}

	time.Sleep(50 * time.Millisecond)

	for _, event := range events {
		r.events = append(r.events, *event)
	}

	log.Printf("MockEventRepository: Successfully inserted %d events. Total stored: %d",
		len(events), len(r.events))

	return nil
}

func (r *MockEventRepository) validateEvent(event *domain.Event) error {
	if event.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if event.MessageID == "" {
		return fmt.Errorf("message_id is required")
	}
	if event.SourceID == "" {
		return fmt.Errorf("source_id is required")
	}
	if event.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if event.EventType == "" {
		return fmt.Errorf("event_type is required")
	}
	if len(event.ContentRaw) == 0 {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}

func (r *MockEventRepository) GetStoredEvents() []domain.Event {
	return r.events
}

func (r *MockEventRepository) Clear() {
	r.events = make([]domain.Event, 0)
}
