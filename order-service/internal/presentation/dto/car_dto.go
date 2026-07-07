package dto

import (
	"time"

	"github.com/qxblvde/CarService/order-service/internal/domain/car"
)

type CreateCarRequest struct {
	CarModelID string `json:"carModelId"`
	VIN        string `json:"vin"`
	Year       int    `json:"year"`
	Price      int64  `json:"price"`
}

type CarResponse struct {
	ID         string    `json:"id"`
	CarModelID string    `json:"carModelId"`
	VIN        string    `json:"vin"`
	Year       int       `json:"year"`
	Price      int64     `json:"price"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func CarToResponse(c car.Car) CarResponse {
	return CarResponse{
		ID:         c.ID,
		CarModelID: c.CarModelID,
		VIN:        c.VIN,
		Year:       c.Year,
		Price:      c.Price.Amount,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func CarsToResponse(cars []car.Car) []CarResponse {
	result := make([]CarResponse, len(cars))
	for i, c := range cars {
		result[i] = CarToResponse(c)
	}
	return result
}
