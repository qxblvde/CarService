package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qxblvde/CarService/storage-service/internal/domain"
)

type AssemblyOrderRepository struct {
	pool *pgxpool.Pool
}

func NewAssemblyOrderRepository(pool *pgxpool.Pool) *AssemblyOrderRepository {
	return &AssemblyOrderRepository{pool: pool}
}

func (r *AssemblyOrderRepository) Create(ctx context.Context, o domain.AssemblyOrder) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO assembly_orders (id, source_order_id, order_type, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		o.ID, o.SourceOrderID, o.OrderType, o.Status, o.CreatedAt, o.UpdatedAt,
	)
	return err
}

func (r *AssemblyOrderRepository) ExistsBySourceOrder(ctx context.Context, sourceOrderID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM assembly_orders WHERE source_order_id = $1)`,
		sourceOrderID).Scan(&exists)
	return exists, err
}

func (r *AssemblyOrderRepository) GetByID(ctx context.Context, id string) (domain.AssemblyOrder, error) {
	var o domain.AssemblyOrder
	err := r.pool.QueryRow(ctx,
		`SELECT id, source_order_id, order_type, status, removed, created_at, updated_at
		 FROM assembly_orders WHERE id = $1`, id).
		Scan(&o.ID, &o.SourceOrderID, &o.OrderType, &o.Status, &o.Removed, &o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.AssemblyOrder{}, pgx.ErrNoRows
	}
	return o, err
}

func (r *AssemblyOrderRepository) List(ctx context.Context) ([]domain.AssemblyOrder, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, source_order_id, order_type, status, removed, created_at, updated_at
		 FROM assembly_orders WHERE removed = false ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.AssemblyOrder
	for rows.Next() {
		var o domain.AssemblyOrder
		if err := rows.Scan(&o.ID, &o.SourceOrderID, &o.OrderType, &o.Status, &o.Removed, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, rows.Err()
}

func (r *AssemblyOrderRepository) UpdateStatus(ctx context.Context, id string, status domain.AssemblyStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE assembly_orders SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	return err
}

func (r *AssemblyOrderRepository) Remove(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE assembly_orders SET removed = true, updated_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}
