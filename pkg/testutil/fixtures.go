package testutil

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/pkg/domain"
)

func NewSampleOrganization(name string) domain.Organization {
	return domain.Organization{
		Identifier: uuid.New(),
		Name:       name,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func NewSampleEvent(tenantID, eventType string) *domain.Event {
	content := map[string]interface{}{
		"user_id": uuid.New().String(),
		"action":  "test_action",
		"value":   42,
	}
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return nil
	}

	return &domain.Event{
		Timestamp:   time.Now(),
		TenantID:    tenantID,
		MessageID:   uuid.New().String(),
		SourceID:    "test-source",
		Topic:       "test-topic",
		ContentRaw:  string(contentJSON),
		ContentJSON: contentJSON,
		ContentSize: uint32(len(contentJSON)),
		Headers: map[string]string{
			"test-header": "test-value",
		},
		FrameType: 1,
		EventType: eventType,
	}
}

func NewBatchEvents(tenantID string, count int, eventType string) []*domain.Event {
	events := make([]*domain.Event, count)
	for i := 0; i < count; i++ {
		events[i] = NewSampleEvent(tenantID, eventType)
	}
	return events
}

func NewBatchIngestionRequest(tenantID string, count int) *domain.BatchIngestionRequest {
	events := make([]domain.EventIngestionData, count)
	for i := 0; i < count; i++ {
		content := map[string]interface{}{
			"index":   i,
			"user_id": uuid.New().String(),
			"action":  fmt.Sprintf("action_%d", i),
		}
		contentJSON, err := json.Marshal(content)
		if err != nil {
			return nil
		}

		events[i] = domain.EventIngestionData{
			MessageID: uuid.New().String(),
			Content:   contentJSON,
			Headers: map[string]string{
				"source": "integration-test",
			},
			FrameType: 1,
			EventType: "test.event",
		}
	}

	return &domain.BatchIngestionRequest{
		TenantID: tenantID,
		SourceID: "integration-test-source",
		Topic:    "integration-test-topic",
		Events:   events,
	}
}
