package persistence

import (
	"context"

	"github.com/qxblvde/CarService/order-service/internal/domain/car"
)

type CarFilter struct {
	CarModelID string
}

type CarRepository interface {
	GetByID(ctx context.Context, id string) (car.Car, error)
	List(ctx context.Context, filter CarFilter) ([]car.Car, error)
	Create(ctx context.Context, c car.Car) error
	Delete(ctx context.Context, id string) error
}
