package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/car"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type carRepository struct {
	pool *pgxpool.Pool
}

func NewCarRepository(pool *pgxpool.Pool) persistence.CarRepository {
	return &carRepository{pool: pool}
}

func (r *carRepository) GetByID(ctx context.Context, id string) (car.Car, error) {
	query := `
		SELECT id, car_model_id, vin, year, price::bigint, created_at, updated_at, deleted_at
		FROM cars
		WHERE id = $1 AND deleted_at IS NULL`

	var c car.Car
	var price int64

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.CarModelID, &c.VIN, &c.Year, &price,
		&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return car.Car{}, apperrs.ErrCarNotFound
		}
		return car.Car{}, err
	}

	c.Price = vo.MustMoney(price)
	return c, nil
}

func (r *carRepository) List(ctx context.Context, filter persistence.CarFilter) ([]car.Car, error) {
	query := `
		SELECT id, car_model_id, vin, year, price::bigint, created_at, updated_at, deleted_at
		FROM cars
		WHERE deleted_at IS NULL
		  AND ($1::uuid IS NULL OR car_model_id = $1::uuid)
		ORDER BY created_at DESC`

	var modelID *string
	if filter.CarModelID != "" {
		modelID = &filter.CarModelID
	}

	rows, err := r.pool.Query(ctx, query, modelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []car.Car
	for rows.Next() {
		var c car.Car
		var price int64
		if err := rows.Scan(
			&c.ID, &c.CarModelID, &c.VIN, &c.Year, &price,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
		); err != nil {
			return nil, err
		}
		c.Price = vo.MustMoney(price)
		result = append(result, c)
	}

	return result, rows.Err()
}

func (r *carRepository) Create(ctx context.Context, c car.Car) error {
	query := `
		INSERT INTO cars (id, car_model_id, vin, year, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		c.ID, c.CarModelID, c.VIN, c.Year, c.Price.Amount,
		c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *carRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE cars
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrs.ErrCarNotFound
	}
	return nil
}
