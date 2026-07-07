package persistence

import (
	"context"

	"github.com/qxblvde/CarService/order-service/internal/domain/user"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (user.User, error)
	Create(ctx context.Context, u user.User) error
}
