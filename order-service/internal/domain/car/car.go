package car

import (
	"time"

	"github.com/qxblvde/CarService/internal/domain/errs"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type Car struct {
	ID         string
	CarModelID string
	VIN        string
	Year       int
	Price      vo.Money
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func NewCar(id, carModelID, vin string, year int, price vo.Money) (Car, error) {
	if id == "" {
		return Car{}, errs.Validation("car id is required")
	}
	if carModelID == "" {
		return Car{}, errs.Validation("car model id is required")
	}
	if vin == "" {
		return Car{}, errs.Validation("vin is required")
	}
	if year <= 1900 {
		return Car{}, errs.Validation("car year must be greater than 1900")
	}
	if price.IsNegative() {
		return Car{}, errs.Validation("car price must be non-negative")
	}

	now := time.Now().UTC()
	return Car{
		ID:         id,
		CarModelID: carModelID,
		VIN:        vin,
		Year:       year,
		Price:      price,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
