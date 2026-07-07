package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/detail"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type detailRepository struct {
	pool *pgxpool.Pool
}

func NewDetailRepository(pool *pgxpool.Pool) persistence.DetailRepository {
	return &detailRepository{pool: pool}
}

func scanDetail(rows pgx.CollectableRow) (detail.Detail, error) {
	var d detail.Detail
	var price int64
	var detailType string
	var compatibleIDs []string

	err := rows.Scan(&d.ID, &d.Name, &detailType, &price, &compatibleIDs)
	if err != nil {
		return detail.Detail{}, err
	}

	d.Price = vo.MustMoney(price)
	d.Type = detail.DetailType(detailType)
	d.CompatibleCarModelIDs = compatibleIDs
	return d, nil
}

const detailQuery = `
	SELECT
	    d.id,
	    d.name,
	    d.type::text,
	    d.price::bigint,
	    COALESCE(
	        ARRAY_AGG(pc.car_model_id::text) FILTER (WHERE pc.car_model_id IS NOT NULL),
	        '{}'
	    ) AS compatible_model_ids
	FROM details d
	LEFT JOIN part_compatibility pc ON d.id = pc.detail_id
	WHERE d.deleted_at IS NULL`

func (r *detailRepository) GetByID(ctx context.Context, id string) (detail.Detail, error) {
	query := detailQuery + ` AND d.id = $1 GROUP BY d.id`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return detail.Detail{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return detail.Detail{}, err
		}
		return detail.Detail{}, apperrs.ErrNotFound
	}

	d, err := scanDetail(rows)
	if err != nil {
		return detail.Detail{}, err
	}
	return d, rows.Err()
}

func (r *detailRepository) List(ctx context.Context) ([]detail.Detail, error) {
	query := detailQuery + ` GROUP BY d.id ORDER BY d.type, d.name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []detail.Detail
	for rows.Next() {
		d, err := scanDetail(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

func (r *detailRepository) ListByModelID(ctx context.Context, carModelID string) ([]detail.Detail, error) {
	query := detailQuery + `
		AND d.id IN (SELECT detail_id FROM part_compatibility WHERE car_model_id = $1)
		GROUP BY d.id
		ORDER BY d.type, d.name`

	rows, err := r.pool.Query(ctx, query, carModelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []detail.Detail
	for rows.Next() {
		d, err := scanDetail(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

func (r *detailRepository) Create(ctx context.Context, d detail.Detail) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO details (id, name, type, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		d.ID, d.Name, string(d.Type), d.Price.Amount, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, modelID := range d.CompatibleCarModelIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO part_compatibility (detail_id, car_model_id) VALUES ($1, $2)`,
			d.ID, modelID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
