package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/test_drive"
)

type TestDriveService struct {
	testDrives persistence.TestDriveRepository
	cars       persistence.CarRepository
}

func NewTestDriveService(testDrives persistence.TestDriveRepository, cars persistence.CarRepository) *TestDriveService {
	return &TestDriveService{testDrives: testDrives, cars: cars}
}

type CreateTestDriveParams struct {
	ClientID      string
	CarID         string
	ScheduledTime time.Time
}

func (s *TestDriveService) Create(ctx context.Context, p CreateTestDriveParams) (test_drive.TestDriveRequest, error) {
	if _, err := s.cars.GetByID(ctx, p.CarID); err != nil {
		return test_drive.TestDriveRequest{}, err
	}

	r, err := test_drive.NewTestDriveRequest(uuid.NewString(), p.CarID, p.ClientID, p.ScheduledTime)
	if err != nil {
		return test_drive.TestDriveRequest{}, err
	}

	if err := s.testDrives.Create(ctx, r); err != nil {
		return test_drive.TestDriveRequest{}, err
	}
	return r, nil
}

func (s *TestDriveService) GetByID(ctx context.Context, id string) (test_drive.TestDriveRequest, error) {
	return s.testDrives.GetByID(ctx, id)
}

func (s *TestDriveService) List(ctx context.Context, filter persistence.TestDriveFilter) ([]test_drive.TestDriveRequest, error) {
	return s.testDrives.List(ctx, filter)
}

func (s *TestDriveService) Cancel(ctx context.Context, id string) error {
	return s.testDrives.Cancel(ctx, id)
}
