package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xtommas/food-backend/internal/data"
)

const (
	outboxIdleDelay   = time.Second
	outboxRetryDelay  = 3 * time.Second
	outboxSendTimeout = 5 * time.Second
)

func (app *application) startOutboxRelay(ctx context.Context) {
	if app.config.rabbitmq.url == "" {
		app.logger.PrintInfo("outbox relay disabled; RABBITMQ_URL is not set", nil)
		return
	}

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		app.logger.PrintInfo("outbox relay started", nil)
		for {
			processed, err := app.models.Outbox.PublishNext(ctx, func(ctx context.Context, event data.OutboxEvent) error {
				publishCtx, cancel := context.WithTimeout(ctx, outboxSendTimeout)
				defer cancel()
				return app.mq.Publish(publishCtx, event.RoutingKey, event.Payload, fmt.Sprintf("%d", event.ID))
			})
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				app.logger.PrintError(err, map[string]string{"component": "outbox_relay"})
				if !waitForContext(ctx, outboxRetryDelay) {
					return
				}
				continue
			}
			if processed {
				continue
			}
			if !waitForContext(ctx, outboxIdleDelay) {
				return
			}
		}
	}()
}

func waitForContext(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
