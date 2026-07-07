package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/qxblvde/CarService/internal/domain/errs"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/service"
)

type OrderResult struct {
	OrderID   string `json:"orderId"`
	OrderType string `json:"orderType"`
	TraceID   string `json:"traceId"`
}

type OrderResultConsumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	svc  *service.OrderService
}

func NewOrderResultConsumer(url string, svc *service.OrderService) (*OrderResultConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &OrderResultConsumer{conn: conn, ch: ch, svc: svc}, nil
}

func (c *OrderResultConsumer) Start(ctx context.Context) error {
	if err := c.consume(ctx, "order.approved", "orders.approved", c.svc.ApplyApproval); err != nil {
		return err
	}
	return c.consume(ctx, "order.rejected", "orders.rejected", c.svc.ApplyRejection)
}

func (c *OrderResultConsumer) consume(ctx context.Context, topic, queue string, apply func(context.Context, string, string) error) error {
	if err := c.ch.ExchangeDeclare(topic, "fanout", true, false, false, false, nil); err != nil {
		return err
	}
	q, err := c.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := c.ch.QueueBind(q.Name, "", topic, false, nil); err != nil {
		return err
	}
	msgs, err := c.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				c.handle(ctx, msg, topic, apply)
			}
		}
	}()
	return nil
}

func (c *OrderResultConsumer) handle(ctx context.Context, msg amqp.Delivery, topic string, apply func(context.Context, string, string) error) {
	var event OrderResult
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("order result decode: %v", err)
		msg.Nack(false, false)
		return
	}

	log.Printf("order result %s: order=%s type=%s trace=%s", topic, event.OrderID, event.OrderType, event.TraceID)

	err := apply(ctx, event.OrderType, event.OrderID)
	switch {
	case err == nil:
		msg.Ack(false)
	case isAlreadyApplied(err):
		log.Printf("order result %s: order=%s already applied (%v)", topic, event.OrderID, err)
		msg.Ack(false)
	default:
		log.Printf("order result %s: order=%s apply failed: %v", topic, event.OrderID, err)
		msg.Nack(false, true)
	}
}

func isAlreadyApplied(err error) bool {
	return errors.Is(err, errs.ErrInvalidTransition) ||
		errors.Is(err, errs.ErrValidation) ||
		errors.Is(err, apperrs.ErrOrderNotFound)
}

func (c *OrderResultConsumer) Close() {
	c.ch.Close()
	c.conn.Close()
}
