package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/detail"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type DetailService struct {
	details persistence.DetailRepository
}

func NewDetailService(details persistence.DetailRepository) *DetailService {
	return &DetailService{details: details}
}

type CreateDetailParams struct {
	Name                  string
	Type                  string
	Price                 int64
	CompatibleCarModelIDs []string
}

func (s *DetailService) Create(ctx context.Context, p CreateDetailParams) (detail.Detail, error) {
	price, err := vo.NewMoney(p.Price)
	if err != nil {
		return detail.Detail{}, err
	}

	d, err := detail.NewDetail(uuid.NewString(), p.Name, detail.DetailType(p.Type), price, p.CompatibleCarModelIDs...)
	if err != nil {
		return detail.Detail{}, err
	}

	if err := s.details.Create(ctx, d); err != nil {
		return detail.Detail{}, err
	}
	return d, nil
}

func (s *DetailService) GetByID(ctx context.Context, id string) (detail.Detail, error) {
	return s.details.GetByID(ctx, id)
}

func (s *DetailService) List(ctx context.Context) ([]detail.Detail, error) {
	return s.details.List(ctx)
}

func (s *DetailService) ListByModel(ctx context.Context, carModelID string) ([]detail.Detail, error) {
	return s.details.ListByModelID(ctx, carModelID)
}
