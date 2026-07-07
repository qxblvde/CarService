package persistence

import (
	"context"

	"github.com/qxblvde/CarService/order-service/internal/domain/car"
)

type CarModelFilter struct {
	Brand     string
	DetailIDs []string
}

type CarModelRepository interface {
	GetByID(ctx context.Context, id string) (car.CarModel, error)
	List(ctx context.Context, filter CarModelFilter) ([]car.CarModel, error)
	Create(ctx context.Context, m car.CarModel) error
	Delete(ctx context.Context, id string) error
}
