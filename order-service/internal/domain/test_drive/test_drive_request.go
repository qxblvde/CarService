package test_drive

import (
	"time"

	"github.com/qxblvde/CarService/internal/domain/errs"
)

type TestDriveRequest struct {
	ID            string
	CarID         string
	ClientID      string
	ScheduledTime time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

func NewTestDriveRequest(id, carID, clientID string, scheduledTime time.Time) (TestDriveRequest, error) {
	if id == "" {
		return TestDriveRequest{}, errs.Validation("test drive request id is required")
	}
	if carID == "" {
		return TestDriveRequest{}, errs.Validation("car id is required")
	}
	if clientID == "" {
		return TestDriveRequest{}, errs.Validation("client id is required")
	}
	if scheduledTime.IsZero() || !scheduledTime.After(time.Now().UTC()) {
		return TestDriveRequest{}, errs.InvalidScheduledTime
	}

	now := time.Now().UTC()
	return TestDriveRequest{
		ID:            id,
		CarID:         carID,
		ClientID:      clientID,
		ScheduledTime: scheduledTime.UTC(),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (r *TestDriveRequest) Reschedule(scheduledTime time.Time) error {
	if scheduledTime.IsZero() || !scheduledTime.After(time.Now().UTC()) {
		return errs.InvalidScheduledTime
	}

	r.ScheduledTime = scheduledTime.UTC()
	r.UpdatedAt = time.Now().UTC()
	return nil
}
