package errors

import "errors"

var (
	ErrNotAClickhouseDriver     = errors.New("not a clickhouse driver")
	ErrEventCouldNotInserted    = errors.New("event could not be inserted")
	ErrEventCouldNotBeValidated = errors.New("event could not be validated")
	ErrCouldNotTransformEvent   = errors.New("event could not be transformed")
	ErrRepositoryCouldNotInsert = errors.New("repository could not be inserted")
)
