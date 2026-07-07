package infrastructure

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxMessage struct {
	ID        string
	Topic     string
	Payload   []byte
	CreatedAt time.Time
}

type OutboxRepository struct {
	pool *pgxpool.Pool
}

func NewOutboxRepository(pool *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{pool: pool}
}

func (r *OutboxRepository) SaveTx(ctx context.Context, tx pgx.Tx, topic string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO outbox (topic, payload) VALUES ($1, $2)`,
		topic, data,
	)
	return err
}

func (r *OutboxRepository) FetchUnpublished(ctx context.Context) ([]OutboxMessage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, topic, payload, created_at
		 FROM outbox WHERE published = false
		 ORDER BY created_at
		 LIMIT 100`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []OutboxMessage
	for rows.Next() {
		var m OutboxMessage
		if err := rows.Scan(&m.ID, &m.Topic, &m.Payload, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

func (r *OutboxRepository) MarkPublished(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE outbox SET published = true WHERE id = $1`,
		id,
	)
	return err
}
