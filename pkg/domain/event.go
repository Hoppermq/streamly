package domain

import (
	"context"
	"encoding/json"
	"time"
)

type Event struct {
	Timestamp     time.Time              `json:"timestamp"`
	TenantID      string                 `json:"tenant_id"`
	MessageID     string                 `json:"message_id"`
	SourceID      string                 `json:"source_id"`
	Topic         string                 `json:"topic"`
	ContentRaw    string                 `json:"-"`
	ContentJSON   json.RawMessage        `json:"content"`
	ContentSize   uint32                 `json:"-"`
	Headers       map[string]string      `json:"headers"`
	FrameType     uint8                  `json:"frame_type"`
	EventType     string                 `json:"event_type"`
}

type BatchIngestionRequest struct {
	TenantID string               `json:"tenant_id" binding:"required"`
	SourceID string               `json:"source_id" binding:"required"`
	Topic    string               `json:"topic" binding:"required"`
	Events   []EventIngestionData `json:"events" binding:"required,min=1,max=5000"`
}

type EventIngestionData struct {
	MessageID string            `json:"message_id" binding:"required"`
	Content   json.RawMessage   `json:"content" binding:"required"`
	Headers   map[string]string `json:"headers"`
	FrameType uint8             `json:"frame_type"`
	EventType string            `json:"event_type" binding:"required"`
}

type BatchIngestionResponse struct {
	Status        string    `json:"status"`
	IngestedCount int       `json:"ingested_count"`
	Timestamp     time.Time `json:"timestamp"`
	BatchID       string    `json:"batch_id"`
	FailedCount   int       `json:"failed_count,omitempty"`
	Failures      []EventFailure `json:"failures,omitempty"`
}

type EventFailure struct {
	EventIndex int    `json:"event_index"`
	MessageID  string `json:"message_id"`
	Error      string `json:"error"`
}

type EventRepository interface {
	BatchInsert(ctx context.Context, events []*Event) error
}

type EventIngestionUseCase interface {
	IngestBatch(ctx context.Context, request *BatchIngestionRequest) (*BatchIngestionResponse, error)
}