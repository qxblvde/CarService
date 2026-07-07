package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qxblvde/CarService/storage-service/internal/domain"
)

type CarRepository struct {
	pool *pgxpool.Pool
}

func NewCarRepository(pool *pgxpool.Pool) *CarRepository {
	return &CarRepository{pool: pool}
}

func (r *CarRepository) ListAvailable(ctx context.Context) ([]domain.Car, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, brand, model, year, price::bigint, available
		 FROM cars WHERE available = TRUE ORDER BY brand, model`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []domain.Car
	for rows.Next() {
		var c domain.Car
		if err := rows.Scan(&c.ID, &c.Brand, &c.Model, &c.Year, &c.Price, &c.Available); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}
	return cars, rows.Err()
}

func (r *CarRepository) GetByID(ctx context.Context, id string) (domain.Car, error) {
	var c domain.Car
	err := r.pool.QueryRow(ctx,
		`SELECT id, brand, model, year, price::bigint, available
		 FROM cars WHERE id = $1 AND available = TRUE`, id).
		Scan(&c.ID, &c.Brand, &c.Model, &c.Year, &c.Price, &c.Available)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Car{}, domain.ErrCarNotFound
	}
	if err != nil {
		return domain.Car{}, err
	}
	return c, nil
}
