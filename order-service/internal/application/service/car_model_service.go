package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/car"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type CarModelService struct {
	models persistence.CarModelRepository
}

func NewCarModelService(models persistence.CarModelRepository) *CarModelService {
	return &CarModelService{models: models}
}

type CreateCarModelParams struct {
	Brand          string
	Model          string
	BasePrice      int64
	BodyType       string
	FuelType       string
	EnginePowerHP  int
	EngineVolumeCC int
	Transmission   string
	DriveType      string
}

func (s *CarModelService) Create(ctx context.Context, p CreateCarModelParams) (car.CarModel, error) {
	price, err := vo.NewMoney(p.BasePrice)
	if err != nil {
		return car.CarModel{}, err
	}

	m, err := car.NewCarModel(
		uuid.NewString(),
		p.Brand, p.Model,
		price,
		car.BodyType(p.BodyType),
		car.FuelType(p.FuelType),
		p.EnginePowerHP,
		p.EngineVolumeCC,
		car.TransmissionType(p.Transmission),
		car.DriveType(p.DriveType),
	)
	if err != nil {
		return car.CarModel{}, err
	}

	if err := s.models.Create(ctx, m); err != nil {
		return car.CarModel{}, err
	}
	return m, nil
}

func (s *CarModelService) GetByID(ctx context.Context, id string) (car.CarModel, error) {
	return s.models.GetByID(ctx, id)
}

func (s *CarModelService) List(ctx context.Context, brand string, detailIDs []string) ([]car.CarModel, error) {
	return s.models.List(ctx, persistence.CarModelFilter{
		Brand:     brand,
		DetailIDs: detailIDs,
	})
}

func (s *CarModelService) Delete(ctx context.Context, id string) error {
	return s.models.Delete(ctx, id)
}
