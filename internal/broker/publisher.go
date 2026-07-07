package broker

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{conn: conn, ch: ch}, nil
}

func (p *Publisher) Publish(ctx context.Context, topic string, body []byte) error {
	if err := p.ch.ExchangeDeclare(topic, "fanout", true, false, false, false, nil); err != nil {
		return err
	}
	return p.ch.PublishWithContext(ctx, topic, "", false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	})
}

func (p *Publisher) Close() {
	if err := p.ch.Close(); err != nil {
		log.Printf("publisher channel close: %v", err)
	}
	if err := p.conn.Close(); err != nil {
		log.Printf("publisher conn close: %v", err)
	}
}
