package messaging

import (
	"context"
	"log"
	"time"

	"github.com/qxblvde/CarService/internal/broker"
	"github.com/qxblvde/CarService/order-service/internal/infrastructure"
)

type OutboxPoller struct {
	repo      *infrastructure.OutboxRepository
	publisher *broker.Publisher
	interval  time.Duration
}

func NewOutboxPoller(repo *infrastructure.OutboxRepository, publisher *broker.Publisher) *OutboxPoller {
	return &OutboxPoller{
		repo:      repo,
		publisher: publisher,
		interval:  3 * time.Second,
	}
}

func (p *OutboxPoller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *OutboxPoller) poll(ctx context.Context) {
	msgs, err := p.repo.FetchUnpublished(ctx)
	if err != nil {
		log.Printf("outbox fetch: %v", err)
		return
	}

	for _, msg := range msgs {
		if err := p.publisher.Publish(ctx, msg.Topic, msg.Payload); err != nil {
			log.Printf("outbox publish %s: %v", msg.ID, err)
			continue
		}
		if err := p.repo.MarkPublished(ctx, msg.ID); err != nil {
			log.Printf("outbox mark published %s: %v", msg.ID, err)
		}
	}
}
