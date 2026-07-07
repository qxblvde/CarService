package service

import (
	"context"

	"github.com/google/uuid"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/car"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type CarService struct {
	cars   persistence.CarRepository
	models persistence.CarModelRepository
}

func NewCarService(cars persistence.CarRepository, models persistence.CarModelRepository) *CarService {
	return &CarService{cars: cars, models: models}
}

type CreateCarParams struct {
	CarModelID string
	VIN        string
	Year       int
	Price      int64
}

func (s *CarService) Create(ctx context.Context, p CreateCarParams) (car.Car, error) {
	if _, err := s.models.GetByID(ctx, p.CarModelID); err != nil {
		return car.Car{}, apperrs.ErrNotFound
	}

	price, err := vo.NewMoney(p.Price)
	if err != nil {
		return car.Car{}, err
	}

	c, err := car.NewCar(uuid.NewString(), p.CarModelID, p.VIN, p.Year, price)
	if err != nil {
		return car.Car{}, err
	}

	if err := s.cars.Create(ctx, c); err != nil {
		return car.Car{}, err
	}
	return c, nil
}

func (s *CarService) GetByID(ctx context.Context, id string) (car.Car, error) {
	return s.cars.GetByID(ctx, id)
}

func (s *CarService) List(ctx context.Context, carModelID string) ([]car.Car, error) {
	return s.cars.List(ctx, persistence.CarFilter{CarModelID: carModelID})
}

func (s *CarService) Delete(ctx context.Context, id string) error {
	return s.cars.Delete(ctx, id)
}
