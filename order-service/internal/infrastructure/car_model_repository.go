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

type carModelRepository struct {
	pool *pgxpool.Pool
}

func NewCarModelRepository(pool *pgxpool.Pool) persistence.CarModelRepository {
	return &carModelRepository{pool: pool}
}

const carModelColumns = `
	id, brand, model, base_price::bigint,
	body_type, fuel_type, engine_power_hp, engine_volume_cc,
	transmission, drive_type,
	created_at, updated_at, deleted_at`

func scanCarModel(row pgx.Row) (car.CarModel, error) {
	var m car.CarModel
	var basePrice int64
	var bodyType, fuelType, transmission, driveType string

	err := row.Scan(
		&m.ID, &m.Brand, &m.Model, &basePrice,
		&bodyType, &fuelType, &m.EnginePowerHP, &m.EngineVolumeCC,
		&transmission, &driveType,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	)
	if err != nil {
		return car.CarModel{}, err
	}

	m.BasePrice = vo.MustMoney(basePrice)
	m.BodyType = car.BodyType(bodyType)
	m.FuelType = car.FuelType(fuelType)
	m.Transmission = car.TransmissionType(transmission)
	m.DriveType = car.DriveType(driveType)
	return m, nil
}

func (r *carModelRepository) GetByID(ctx context.Context, id string) (car.CarModel, error) {
	query := `SELECT ` + carModelColumns + `
		FROM car_models
		WHERE id = $1 AND deleted_at IS NULL`

	m, err := scanCarModel(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return car.CarModel{}, apperrs.ErrNotFound
		}
		return car.CarModel{}, err
	}
	return m, nil
}

func (r *carModelRepository) List(ctx context.Context, filter persistence.CarModelFilter) ([]car.CarModel, error) {
	query := `SELECT ` + carModelColumns + `
		FROM car_models cm
		WHERE cm.deleted_at IS NULL
		  AND ($1 = '' OR cm.brand = $1)
		  AND ($2::uuid[] IS NULL OR cm.id IN (
		        SELECT car_model_id FROM part_compatibility WHERE detail_id = ANY($2::uuid[])
		      ))
		ORDER BY cm.brand, cm.model`

	var detailIDs interface{} = nil
	if len(filter.DetailIDs) > 0 {
		detailIDs = filter.DetailIDs
	}

	rows, err := r.pool.Query(ctx, query, filter.Brand, detailIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []car.CarModel
	for rows.Next() {
		m, err := scanCarModel(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func (r *carModelRepository) Create(ctx context.Context, m car.CarModel) error {
	query := `
		INSERT INTO car_models
		    (id, brand, model, base_price, body_type, fuel_type,
		     engine_power_hp, engine_volume_cc, transmission, drive_type,
		     created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	_, err := r.pool.Exec(ctx, query,
		m.ID, m.Brand, m.Model, m.BasePrice.Amount,
		string(m.BodyType), string(m.FuelType),
		m.EnginePowerHP, m.EngineVolumeCC,
		string(m.Transmission), string(m.DriveType),
		m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *carModelRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE car_models
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrs.ErrNotFound
	}
	return nil
}
