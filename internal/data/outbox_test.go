package data

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestOrderModel_InsertWithEvent(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := newTestOrder(user.Id, restaurantID)

	var payload string
	err := orderModel.InsertWithEvent(order, "order.created", func() ([]byte, error) {
		payload = fmt.Sprintf(`{"order_id": %d}`, order.ID)
		return []byte(payload), nil
	})
	if err != nil {
		t.Fatalf("InsertWithEvent() error = %v", err)
	}
	var outboxID int64
	var routingKey string
	if err := testDB.QueryRow(`
		SELECT id, routing_key, payload::text
		FROM outbox_events
		WHERE payload::text = $1`, payload).Scan(&outboxID, &routingKey, &payload); err != nil {
		t.Fatalf("fetching outbox event: %v", err)
	}
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM outbox_events WHERE id = $1`, outboxID)
		testDB.Exec(`DELETE FROM orders WHERE id = $1`, order.ID)
	})
	if routingKey != "order.created" {
		t.Errorf("routing_key = %q, want order.created", routingKey)
	}
	if payload != fmt.Sprintf(`{"order_id": %d}`, order.ID) {
		t.Errorf("payload = %q, want order JSON", payload)
	}
}

func TestOutboxModel_PublishNext(t *testing.T) {
	var eventID int64
	err := testDB.QueryRow(`
		INSERT INTO outbox_events (routing_key, payload)
		VALUES ('order.created', '{"order_id": 123456789}'::jsonb)
		RETURNING id`).Scan(&eventID)
	if err != nil {
		t.Fatalf("creating outbox event: %v", err)
	}
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM outbox_events WHERE id = $1`, eventID)
	})

	model := OutboxModel{DB: testDB}
	processed, err := model.PublishNext(context.Background(), func(_ context.Context, event OutboxEvent) error {
		if event.RoutingKey != "order.created" {
			t.Errorf("RoutingKey = %q, want order.created", event.RoutingKey)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("PublishNext() error = %v", err)
	}
	if !processed {
		t.Fatal("PublishNext() processed = false, want true")
	}

	var published bool
	var attempts int
	if err := testDB.QueryRow(`
		SELECT published_at IS NOT NULL, attempts
		FROM outbox_events
		WHERE id = $1`, eventID).Scan(&published, &attempts); err != nil {
		t.Fatalf("fetching published event: %v", err)
	}
	if !published || attempts != 1 {
		t.Errorf("published = %t, attempts = %d, want true and 1", published, attempts)
	}
}

func TestOutboxModel_PublishNext_RetainsFailedEvent(t *testing.T) {
	var eventID int64
	err := testDB.QueryRow(`
		INSERT INTO outbox_events (routing_key, payload)
		VALUES ('order.created', '{"order_id": 987654321}'::jsonb)
		RETURNING id`).Scan(&eventID)
	if err != nil {
		t.Fatalf("creating outbox event: %v", err)
	}
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM outbox_events WHERE id = $1`, eventID)
	})

	publishErr := errors.New("rabbitmq unavailable")
	model := OutboxModel{DB: testDB}
	processed, err := model.PublishNext(context.Background(), func(context.Context, OutboxEvent) error {
		return publishErr
	})
	if !processed {
		t.Fatal("PublishNext() processed = false, want true")
	}
	if !errors.Is(err, publishErr) {
		t.Fatalf("PublishNext() error = %v, want %v", err, publishErr)
	}

	var published bool
	var attempts int
	var lastError string
	if err := testDB.QueryRow(`
		SELECT published_at IS NOT NULL, attempts, last_error
		FROM outbox_events
		WHERE id = $1`, eventID).Scan(&published, &attempts, &lastError); err != nil {
		t.Fatalf("fetching failed event: %v", err)
	}
	if published || attempts != 1 || lastError != publishErr.Error() {
		t.Errorf("published = %t, attempts = %d, last_error = %q", published, attempts, lastError)
	}
}
