package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/order"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type inStockOrderRepository struct {
	pool *pgxpool.Pool
}

func NewInStockOrderRepository(pool *pgxpool.Pool) persistence.InStockOrderRepository {
	return &inStockOrderRepository{pool: pool}
}

func scanInStockOrder(row pgx.Row) (order.InStockOrder, error) {
	var o order.InStockOrder
	var price int64
	var status string

	err := row.Scan(
		&o.ID, &o.ManagerID, &o.ClientID, &o.CarID,
		&status, &price,
		&o.CreatedAt, &o.UpdatedAt, &o.DeletedAt,
	)
	if err != nil {
		return order.InStockOrder{}, err
	}

	o.Price = vo.MustMoney(price)
	o.Status = order.InStockStatus(status)
	return o, nil
}

const inStockOrderSelect = `
	SELECT id, manager_id, client_id, car_id,
	       status::text, price::bigint,
	       created_at, updated_at, deleted_at
	FROM in_stock_orders
	WHERE deleted_at IS NULL`

func (r *inStockOrderRepository) GetByID(ctx context.Context, id string) (order.InStockOrder, error) {
	row := r.pool.QueryRow(ctx, inStockOrderSelect+` AND id = $1`, id)
	o, err := scanInStockOrder(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order.InStockOrder{}, apperrs.ErrOrderNotFound
		}
		return order.InStockOrder{}, err
	}
	return o, nil
}

func (r *inStockOrderRepository) List(ctx context.Context, filter persistence.OrderFilter) ([]order.InStockOrder, error) {
	query := inStockOrderSelect + `
		AND ($1 = '' OR client_id = $1::uuid)
		AND ($2 = '' OR manager_id = $2::uuid)
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, filter.ClientID, filter.ManagerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []order.InStockOrder
	for rows.Next() {
		o, err := scanInStockOrder(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, rows.Err()
}

func (r *inStockOrderRepository) Create(ctx context.Context, o order.InStockOrder) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO in_stock_orders
		    (id, manager_id, client_id, car_id, status, price, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		o.ID, o.ManagerID, o.ClientID, o.CarID,
		string(o.Status), o.Price.Amount,
		o.CreatedAt, o.UpdatedAt,
	)
	return err
}

func (r *inStockOrderRepository) UpdateStatus(ctx context.Context, id string, status order.InStockStatus) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE in_stock_orders
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`,
		id, string(status),
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrs.ErrOrderNotFound
	}
	return nil
}

func (r *inStockOrderRepository) UpdateStatusTx(ctx context.Context, tx Tx, id string, status order.InStockStatus) error {
	tag, err := tx.Exec(ctx, `
		UPDATE in_stock_orders
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`,
		id, string(status),
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrs.ErrOrderNotFound
	}
	return nil
}
