package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/test_drive"
)

type testDriveRepository struct {
	pool *pgxpool.Pool
}

func NewTestDriveRepository(pool *pgxpool.Pool) persistence.TestDriveRepository {
	return &testDriveRepository{pool: pool}
}

const testDriveSelect = `
	SELECT id, client_id, car_id, scheduled_time, created_at, updated_at, deleted_at
	FROM test_drive_requests
	WHERE deleted_at IS NULL`

func scanTestDrive(row pgx.Row) (test_drive.TestDriveRequest, error) {
	var r test_drive.TestDriveRequest
	return r, row.Scan(
		&r.ID, &r.ClientID, &r.CarID, &r.ScheduledTime,
		&r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
	)
}

func (r *testDriveRepository) GetByID(ctx context.Context, id string) (test_drive.TestDriveRequest, error) {
	td, err := scanTestDrive(r.pool.QueryRow(ctx, testDriveSelect+` AND id = $1`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return test_drive.TestDriveRequest{}, apperrs.ErrTestDriveNotFound
		}
		return test_drive.TestDriveRequest{}, err
	}
	return td, nil
}

func (r *testDriveRepository) List(ctx context.Context, filter persistence.TestDriveFilter) ([]test_drive.TestDriveRequest, error) {
	query := testDriveSelect + `
		AND ($1 = '' OR client_id = $1::uuid)
		AND ($2 = '' OR car_id = $2::uuid)
		ORDER BY scheduled_time`

	rows, err := r.pool.Query(ctx, query, filter.ClientID, filter.CarID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []test_drive.TestDriveRequest
	for rows.Next() {
		td, err := scanTestDrive(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, td)
	}
	return result, rows.Err()
}

func (r *testDriveRepository) Create(ctx context.Context, td test_drive.TestDriveRequest) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO test_drive_requests
		    (id, client_id, car_id, scheduled_time, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		td.ID, td.ClientID, td.CarID, td.ScheduledTime,
		td.CreatedAt, td.UpdatedAt,
	)
	return err
}

func (r *testDriveRepository) Cancel(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE test_drive_requests
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrs.ErrTestDriveNotFound
	}
	return nil
}
