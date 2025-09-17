package errors

import "errors"

var (
	ErrEngineErrorOrder         = errors.New("engine should be set before server")
	ErrEmptyContent             = errors.New("content could not be empty")
	ErrBatchSizeMaxSizeExceeded = errors.New("batch size exceeded")
)
