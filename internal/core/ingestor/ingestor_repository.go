package ingestor

import (
	"context"

	"github.com/hoppermq/streamly/internal/storage/clickhouse"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/errors"
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
	if len(events) == 0 {
		return nil
	}

	// TODO: add proper methods to the domain.
	driver, ok := e.driver.(*clickhouse.Driver)
	if !ok {
		return errors.ErrNotAClickhouseDriver
	}

	tx, err := driver.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	committed := false
	defer func(tx domain.Tx) {
		if !committed {
			if err = tx.Rollback(); err != nil {
				return
			}
		}
	}(tx)

	// We will extract query to another place.
	query := `INSERT INTO events (
          timestamp, tenant_id, message_id, source_id, topic,
          content_raw, content_json, content_size_bytes, headers, frame_type, event_type
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func(stmt domain.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	for _, event := range events {
		err := stmt.ExecContext(ctx,
			event.Timestamp,
			event.TenantID,
			event.MessageID,
			event.SourceID,
			event.Topic,
			event.ContentRaw,
			string(event.ContentJSON),
			event.ContentSize,
			event.Headers,
			event.FrameType,
			event.EventType,
		)
		if err != nil {
			return errors.ErrEventCouldNotInserted
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	committed = true
	return nil
}

func NewEventRepository(opts ...RepositoryOption) *EventRepository {
	r := &EventRepository{}
	for _, opt := range opts {
		opt(r)
	}

	return r
}
