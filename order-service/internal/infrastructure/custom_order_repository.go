package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/car"
	"github.com/qxblvde/CarService/order-service/internal/domain/detail"
	"github.com/qxblvde/CarService/order-service/internal/domain/order"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type customOrderRepository struct {
	pool *pgxpool.Pool
}

func NewCustomOrderRepository(pool *pgxpool.Pool) persistence.CustomOrderRepository {
	return &customOrderRepository{pool: pool}
}

func (r *customOrderRepository) loadParts(ctx context.Context, orderIDs []string) (map[string][]detail.Detail, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT cop.order_id, d.id, d.name, d.type::text, d.price::bigint
		FROM custom_order_parts cop
		JOIN details d ON d.id = cop.detail_id
		WHERE cop.order_id = ANY($1::uuid[])`,
		orderIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]detail.Detail)
	for rows.Next() {
		var orderID string
		var d detail.Detail
		var price int64
		var detailType string

		if err := rows.Scan(&orderID, &d.ID, &d.Name, &detailType, &price); err != nil {
			return nil, err
		}
		d.Type = detail.DetailType(detailType)
		d.Price = vo.MustMoney(price)
		result[orderID] = append(result[orderID], d)
	}
	return result, rows.Err()
}

func buildConfig(carModelID string, parts []detail.Detail) car.Configuration {
	cfg := car.Configuration{
		CarModelID: carModelID,
		Parts:      make(map[detail.DetailType]detail.Detail, len(parts)),
	}
	for _, p := range parts {
		cfg.Parts[p.Type] = p
	}
	return cfg
}

const customOrderSelect = `
	SELECT id, manager_id, client_id, car_model_id,
	       status::text, price::bigint,
	       created_at, updated_at, deleted_at
	FROM custom_orders
	WHERE deleted_at IS NULL`

func scanCustomOrderRow(row pgx.Row) (order.CustomOrder, error) {
	var o order.CustomOrder
	var price int64
	var status string

	err := row.Scan(
		&o.ID, &o.ManagerID, &o.ClientID, &o.CarModelID,
		&status, &price,
		&o.CreatedAt, &o.UpdatedAt, &o.DeletedAt,
	)
	if err != nil {
		return order.CustomOrder{}, err
	}

	o.Price = vo.MustMoney(price)
	o.Status = order.CustomStatus(status)
	return o, nil
}

func (r *customOrderRepository) GetByID(ctx context.Context, id string) (order.CustomOrder, error) {
	o, err := scanCustomOrderRow(r.pool.QueryRow(ctx, customOrderSelect+` AND id = $1`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order.CustomOrder{}, apperrs.ErrOrderNotFound
		}
		return order.CustomOrder{}, err
	}

	parts, err := r.loadParts(ctx, []string{o.ID})
	if err != nil {
		return order.CustomOrder{}, err
	}
	o.Configuration = buildConfig(o.CarModelID, parts[o.ID])
	return o, nil
}

func (r *customOrderRepository) List(ctx context.Context, filter persistence.OrderFilter) ([]order.CustomOrder, error) {
	query := customOrderSelect + `
		AND ($1 = '' OR client_id = $1::uuid)
		AND ($2 = '' OR manager_id = $2::uuid)
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, filter.ClientID, filter.ManagerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []order.CustomOrder
	var ids []string
	for rows.Next() {
		o, err := scanCustomOrderRow(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
		ids = append(ids, o.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return orders, nil
	}

	partsByOrder, err := r.loadParts(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range orders {
		orders[i].Configuration = buildConfig(orders[i].CarModelID, partsByOrder[orders[i].ID])
	}
	return orders, nil
}

func (r *customOrderRepository) Create(ctx context.Context, o order.CustomOrder) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO custom_orders
		    (id, manager_id, client_id, car_model_id, status, price, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		o.ID, o.ManagerID, o.ClientID, o.CarModelID,
		string(o.Status), o.Price.Amount,
		o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, part := range o.Configuration.Parts {
		_, err = tx.Exec(ctx,
			`INSERT INTO custom_order_parts (order_id, detail_id) VALUES ($1, $2)`,
			o.ID, part.ID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *customOrderRepository) UpdateStatus(ctx context.Context, id string, status order.CustomStatus) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE custom_orders
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

func (r *customOrderRepository) UpdateStatusTx(ctx context.Context, tx Tx, id string, status order.CustomStatus) error {
	tag, err := tx.Exec(ctx, `
		UPDATE custom_orders
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
