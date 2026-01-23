package errors

import "errors"

var (
	ErrEngineErrorOrder         = errors.New("engine should be set before server")
	ErrEmptyContent             = errors.New("content could not be empty")
	ErrBatchSizeMaxSizeExceeded = errors.New("batch size exceeded")
)

var (
	ErrTenantIDRequired   = errors.New("tenant_id is required")
	ErrSourceIDRequired   = errors.New("source_id is required")
	ErrMessageIDRequired  = errors.New("message_id is required")
	ErrTopicRequired      = errors.New("topic is required")
	ErrEventTypeRequired  = errors.New("event_type is required")
	ErrRawContentRequired = errors.New("raw_content is required")
	ErrEventEmpty         = errors.New("event cannot be empty")
)
