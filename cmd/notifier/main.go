package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("connecting to rabbitmq: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("opening channel: %v", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare("orders", "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("declaring exchange: %v", err)
	}

	q, err := ch.QueueDeclare("notifier.orders", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("declaring queue: %v", err)
	}

	if err := ch.QueueBind(q.Name, "order.*", "orders", false, nil); err != nil {
		log.Fatalf("binding queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("registering consumer: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	log.Println("notifier started, waiting for messages...")

	go func() {
		for d := range msgs {
			handleMessage(d)
		}
	}()

	<-stop
	log.Println("shutting down notifier")
}

func handleMessage(d amqp.Delivery) {
	switch d.RoutingKey {
	case "order.created":
		var evt struct {
			OrderID      int64  `json:"order_id"`
			RestaurantID int64  `json:"restaurant_id"`
			Address      string `json:"address"`
		}
		if err := json.Unmarshal(d.Body, &evt); err != nil {
			log.Printf("failed to decode order.created: %v", err)
			d.Nack(false, false)
			return
		}
		log.Printf("[kitchen] new order #%d for restaurant %d -> deliver to %q",
			evt.OrderID, evt.RestaurantID, evt.Address)

	case "order.status_changed":
		var evt struct {
			OrderID    int64  `json:"order_id"`
			UserID     int64  `json:"user_id"`
			FromStatus string `json:"from_status"`
			ToStatus   string `json:"to_status"`
		}
		if err := json.Unmarshal(d.Body, &evt); err != nil {
			log.Printf("failed to decode order.status_changed: %v", err)
			d.Nack(false, false)
			return
		}
		log.Printf("[customer] order #%d for user %d: %s -> %s",
			evt.OrderID, evt.UserID, evt.FromStatus, evt.ToStatus)

	default:
		log.Printf("unhandled routing key: %s", d.RoutingKey)
	}

	d.Ack(false)
}
