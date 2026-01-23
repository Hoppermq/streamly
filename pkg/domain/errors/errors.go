package errors

import (
	"errors"
	"fmt"
)

var (
	ErrEngineErrorOrder         = errors.New("engine should be set before server")
	ErrEmptyContent             = errors.New("content could not be empty")
	ErrBatchSizeMaxSizeExceeded = errors.New("batch size exceeded")
	ErrFailedToReadFile         = errors.New("failed to read file")
)

func FailedToReadFile(path string) error {
	return fmt.Errorf("%w: %s", ErrFailedToReadFile, path)
}

var (
	ErrTenantIDRequired   = errors.New("tenant_id is required")
	ErrSourceIDRequired   = errors.New("source_id is required")
	ErrMessageIDRequired  = errors.New("message_id is required")
	ErrTopicRequired      = errors.New("topic is required")
	ErrEventTypeRequired  = errors.New("event_type is required")
	ErrRawContentRequired = errors.New("raw_content is required")
	ErrEventEmpty         = errors.New("event cannot be empty")
)

func EventMessageMissing(eventID int) error {
	return fmt.Errorf("%w: %d", ErrMessageIDRequired, eventID)
}

func EventTypeMissing(eventID int) error {
	return fmt.Errorf("%w: %d", ErrEventTypeRequired, eventID)
}

func EventContentEmpty(eventID int) error {
	return fmt.Errorf("%w: %d", ErrEventEmpty, eventID)
}

var (
	ErrZitadelClientCreation = errors.New("failed to create Zitadel client")
	ErrZitadelPATRequired    = errors.New("pat token is required (use WithPAT or WithPATFromFile)")
)

var (
	ErrSerializerInvalidTimeWindow    = errors.New("invalid time window")
	ErrSerializerInvalidGroupByClause = errors.New("invalid group by clause")
	ErrSerializerInvalidGroupBy       = errors.New("groupBy must be string or time window object")
	ErrSerializerInvalidSelectClause  = errors.New("invalid select clause")
	ErrSerializerInvalidSelect        = errors.New("select must be string or aggregation object")
)

func SerializerInvalidTimeWindow(err error) error {
	return fmt.Errorf("%w: %s", ErrSerializerInvalidTimeWindow, err)
}

func SerializerInvalidSelectFunction(err error) error {
	return fmt.Errorf("%w: %s", ErrSerializerInvalidSelect, err)
}
