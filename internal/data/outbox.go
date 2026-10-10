package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// OutboxEvent is a durable event awaiting delivery to the message broker.
// Consumers should use ID as an idempotency key because delivery is at-least-once.
type OutboxEvent struct {
	ID         int64
	RoutingKey string
	Payload    []byte
	CreatedAt  time.Time
}

type OutboxPublisher func(context.Context, OutboxEvent) error

type OutboxModel struct {
	DB *sql.DB
}

func insertOutboxEvent(ctx context.Context, tx *sql.Tx, routingKey string, payload []byte) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO outbox_events (routing_key, payload)
		VALUES ($1, $2::jsonb)`, routingKey, payload)
	return err
}

// PublishNext locks one event until its successful delivery or failed attempt is
// recorded. A broker acknowledgement followed by a database failure can produce a
// duplicate delivery, so consumers must use ID as an idempotency key.
func (m OutboxModel) PublishNext(ctx context.Context, publish OutboxPublisher) (bool, error) {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var event OutboxEvent
	err = tx.QueryRowContext(ctx, `
		SELECT id, routing_key, payload, created_at
		FROM outbox_events
		WHERE published_at IS NULL
		ORDER BY id
		FOR UPDATE SKIP LOCKED
		LIMIT 1`).Scan(&event.ID, &event.RoutingKey, &event.Payload, &event.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}

	if err := publish(ctx, event); err != nil {
		if _, updateErr := tx.ExecContext(ctx, `
			UPDATE outbox_events
			SET attempts = attempts + 1, last_error = $2
			WHERE id = $1`, event.ID, err.Error()); updateErr != nil {
			return true, fmt.Errorf("recording failed outbox event: %w", updateErr)
		}
		if commitErr := tx.Commit(); commitErr != nil {
			return true, commitErr
		}
		return true, fmt.Errorf("publishing outbox event %d: %w", event.ID, err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE outbox_events
		SET published_at = NOW(), attempts = attempts + 1, last_error = NULL
		WHERE id = $1`, event.ID); err != nil {
		return true, err
	}

	return true, tx.Commit()
}
