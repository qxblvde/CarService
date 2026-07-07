package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/qxblvde/CarService/storage-service/internal/domain"
	"github.com/qxblvde/CarService/storage-service/internal/infrastructure"
)

type AssemblyOrderService struct {
	repo *infrastructure.AssemblyOrderRepository
}

func NewAssemblyOrderService(repo *infrastructure.AssemblyOrderRepository) *AssemblyOrderService {
	return &AssemblyOrderService{repo: repo}
}

func (s *AssemblyOrderService) Create(ctx context.Context, sourceOrderID string, orderType domain.SourceOrderType) (domain.AssemblyOrder, error) {
	o, err := domain.NewAssemblyOrder(uuid.NewString(), sourceOrderID, orderType)
	if err != nil {
		return domain.AssemblyOrder{}, err
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return domain.AssemblyOrder{}, err
	}
	return o, nil
}

func (s *AssemblyOrderService) ExistsBySourceOrder(ctx context.Context, sourceOrderID string) (bool, error) {
	return s.repo.ExistsBySourceOrder(ctx, sourceOrderID)
}

func (s *AssemblyOrderService) GetByID(ctx context.Context, id string) (domain.AssemblyOrder, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AssemblyOrderService) List(ctx context.Context) ([]domain.AssemblyOrder, error) {
	return s.repo.List(ctx)
}

func (s *AssemblyOrderService) Assemble(ctx context.Context, id string) error {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := o.Assemble(); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, o.Status)
}

func (s *AssemblyOrderService) Fail(ctx context.Context, id string) error {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := o.Fail(); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, o.Status)
}

func (s *AssemblyOrderService) Remove(ctx context.Context, id string) error {
	return s.repo.Remove(ctx, id)
}
