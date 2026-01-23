package domain

import (
	"context"
	"encoding/json"
	"time"
)

// Event type represent the ingested event from a broker.
type Event struct {
	Timestamp   time.Time         `json:"timestamp"`
	TenantID    string            `json:"tenant_id"`
	MessageID   string            `json:"message_id"`
	SourceID    string            `json:"source_id"`
	Topic       string            `json:"topic"`
	ContentRaw  string            `json:"-"`
	ContentJSON json.RawMessage   `json:"content"`
	ContentSize uint32            `json:"-"`
	Headers     map[string]string `json:"headers"`
	FrameType   uint8             `json:"frame_type"`
	EventType   string            `json:"event_type"`
}

// BatchIngestionRequest type represent the request resource for ingesting batch of events.
type BatchIngestionRequest struct {
	TenantID string               `binding:"required"                json:"tenant_id"`
	SourceID string               `binding:"required"                json:"source_id"`
	Topic    string               `binding:"required"                json:"topic"`
	Events   []EventIngestionData `binding:"required,min=1,max=5000" json:"events"`
}

// EventIngestionData type represent the ingested message content data.
type EventIngestionData struct {
	MessageID string            `binding:"required" json:"message_id"`
	Content   json.RawMessage   `binding:"required" json:"content"`
	Headers   map[string]string `json:"headers"`
	FrameType uint8             `json:"frame_type"`
	EventType string            `binding:"required" json:"event_type"`
}

// BatchIngestionResponse type represent the response returned after ingesting a batch of events.
type BatchIngestionResponse struct {
	Status        string         `json:"status"`
	IngestedCount int            `json:"ingested_count"`
	Timestamp     time.Time      `json:"timestamp"`
	BatchID       string         `json:"batch_id"`
	FailedCount   int            `json:"failed_count,omitempty"`
	Failures      []EventFailure `json:"failures,omitempty"`
}

// EventFailure type represent the failure of event ingested.
type EventFailure struct {
	EventIndex int    `json:"event_index"`
	MessageID  string `json:"message_id"`
	Error      string `json:"error"`
}

// IngestionRepository type represent the contract for ingestion repository.
type IngestionRepository interface {
	BatchInsert(ctx context.Context, events []*Event) error
}

// IngestionUseCase type represent the contract for ingestion use case.
type IngestionUseCase interface { // should mre be IngestionUseCase
	IngestBatch(ctx context.Context, request *BatchIngestionRequest) (*BatchIngestionResponse, error)
}
