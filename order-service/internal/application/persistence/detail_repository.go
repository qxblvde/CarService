package persistence

import (
	"context"

	"github.com/qxblvde/CarService/order-service/internal/domain/detail"
)

type DetailRepository interface {
	GetByID(ctx context.Context, id string) (detail.Detail, error)
	List(ctx context.Context) ([]detail.Detail, error)
	ListByModelID(ctx context.Context, carModelID string) ([]detail.Detail, error)
	Create(ctx context.Context, d detail.Detail) error
}
