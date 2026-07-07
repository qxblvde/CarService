package messaging

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/qxblvde/CarService/internal/broker"
	"github.com/qxblvde/CarService/storage-service/internal/application/service"
	"github.com/qxblvde/CarService/storage-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderSentForApproval struct {
	OrderID    string    `json:"orderId"`
	OrderType  string    `json:"orderType"`
	TraceID    string    `json:"traceId"`
	OccurredAt time.Time `json:"occurredAt"`
}

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	svc  *service.AssemblyOrderService
	pub  *broker.Publisher
}

func NewConsumer(url string, svc *service.AssemblyOrderService, pub *broker.Publisher) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{conn: conn, ch: ch, svc: svc, pub: pub}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	topic := "order.sent_for_approval"
	queue := "storage.order_approval"

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
				c.handle(ctx, msg)
			}
		}
	}()

	return nil
}

func (c *Consumer) handle(ctx context.Context, msg amqp.Delivery) {
	var event OrderSentForApproval
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("consumer decode: %v", err)
		msg.Nack(false, false)
		return
	}

	if exists, err := c.svc.ExistsBySourceOrder(ctx, event.OrderID); err != nil {
		log.Printf("consumer exists check: %v", err)
		msg.Nack(false, true)
		return
	} else if exists {
		log.Printf("consumer: assembly for order %s already exists, skipping", event.OrderID)
		msg.Ack(false)
		return
	}

	orderType := domain.OrderTypeInStock
	if event.OrderType == "CUSTOM" {
		orderType = domain.OrderTypeCustom
	}

	assembly, err := c.svc.Create(ctx, event.OrderID, orderType)
	if err != nil {
		log.Printf("consumer create assembly: %v", err)
		msg.Nack(false, true)
		return
	}

	if err := c.svc.Assemble(ctx, assembly.ID); err != nil {
		log.Printf("consumer assemble: %v", err)
		_ = c.svc.Fail(ctx, assembly.ID)
		c.publishResult(ctx, event, "order.rejected")
		msg.Ack(false)
		return
	}

	c.publishResult(ctx, event, "order.approved")
	msg.Ack(false)
}

func (c *Consumer) publishResult(ctx context.Context, event OrderSentForApproval, topic string) {
	payload, _ := json.Marshal(map[string]string{
		"orderId":   event.OrderID,
		"orderType": event.OrderType,
		"traceId":   event.TraceID,
	})
	if err := c.pub.Publish(ctx, topic, payload); err != nil {
		log.Printf("consumer publish result: %v", err)
	}
}

func (c *Consumer) Close() {
	c.ch.Close()
	c.conn.Close()
}
