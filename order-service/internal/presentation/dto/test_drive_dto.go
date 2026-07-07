package dto

import (
	"time"

	"github.com/qxblvde/CarService/order-service/internal/domain/test_drive"
)

type CreateTestDriveRequest struct {
	ClientID      string    `json:"clientId"`
	CarID         string    `json:"carId"`
	ScheduledTime time.Time `json:"scheduledTime"`
}

type TestDriveResponse struct {
	ID            string    `json:"id"`
	ClientID      string    `json:"clientId"`
	CarID         string    `json:"carId"`
	ScheduledTime time.Time `json:"scheduledTime"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func TestDriveToResponse(r test_drive.TestDriveRequest) TestDriveResponse {
	return TestDriveResponse{
		ID:            r.ID,
		ClientID:      r.ClientID,
		CarID:         r.CarID,
		ScheduledTime: r.ScheduledTime,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func TestDrivesToResponse(reqs []test_drive.TestDriveRequest) []TestDriveResponse {
	result := make([]TestDriveResponse, len(reqs))
	for i, r := range reqs {
		result[i] = TestDriveToResponse(r)
	}
	return result
}
