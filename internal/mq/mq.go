package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ExchangeName = "orders"

type Publisher struct {
	url string

	mu       sync.Mutex
	conn     *amqp.Connection
	ch       *amqp.Channel
	confirms <-chan amqp.Confirmation
	returns  <-chan amqp.Return
}

func NewPublisher(url string) *Publisher {
	return &Publisher{url: url}
}

func (p *Publisher) connect() error {
	if p.conn != nil && p.ch != nil {
		return nil
	}

	conn, err := amqp.Dial(p.url)
	if err != nil {
		return fmt.Errorf("connecting to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("opening channel: %w", err)
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("declaring exchange: %w", err)
	}
	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("enabling publisher confirms: %w", err)
	}

	p.conn = conn
	p.ch = ch
	p.confirms = ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	p.returns = ch.NotifyReturn(make(chan amqp.Return, 1))
	return nil
}

func (p *Publisher) disconnect() {
	if p.ch != nil {
		_ = p.ch.Close()
		p.ch = nil
		p.confirms = nil
		p.returns = nil
	}
	if p.conn != nil {
		_ = p.conn.Close()
		p.conn = nil
	}
}

func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.disconnect()
	return nil
}

// Publish waits for the broker's publisher confirmation. Calls are serialized
// because confirms belong to an AMQP channel rather than an individual request.
func (p *Publisher) Publish(ctx context.Context, routingKey string, body []byte, messageID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.connect(); err != nil {
		return err
	}

	select {
	case <-p.confirms:
	default:
	}
	select {
	case <-p.returns:
	default:
	}

	err := p.ch.PublishWithContext(ctx,
		ExchangeName,
		routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now().UTC(),
			MessageId:    messageID,
			Body:         body,
		},
	)
	if err != nil {
		p.disconnect()
		return fmt.Errorf("publishing message: %w", err)
	}

	select {
	case confirmation, ok := <-p.confirms:
		if !ok || !confirmation.Ack {
			p.disconnect()
			return fmt.Errorf("broker did not confirm message")
		}
		select {
		case returned, ok := <-p.returns:
			if !ok {
				p.disconnect()
				return fmt.Errorf("broker return channel closed")
			}
			return fmt.Errorf("message was unroutable: %s", returned.ReplyText)
		default:
		}
		return nil
	case <-ctx.Done():
		p.disconnect()
		return ctx.Err()
	}
}
