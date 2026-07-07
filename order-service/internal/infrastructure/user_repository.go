package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/user"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) persistence.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (user.User, error) {
	var u user.User
	var role string

	err := r.pool.QueryRow(ctx, `
		SELECT id, role::text, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&u.ID, &role, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, apperrs.ErrUserNotFound
		}
		return user.User{}, err
	}

	u.Role = user.UserRole(role)
	return u, nil
}

func (r *userRepository) Create(ctx context.Context, u user.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`,
		u.ID, string(u.Role), u.CreatedAt, u.UpdatedAt,
	)
	return err
}
