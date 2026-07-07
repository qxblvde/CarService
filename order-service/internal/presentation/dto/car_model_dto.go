package dto

import (
	"time"

	"github.com/qxblvde/CarService/order-service/internal/domain/car"
)

type CreateCarModelRequest struct {
	Brand          string `json:"brand"`
	Model          string `json:"model"`
	BasePrice      int64  `json:"basePrice"`
	BodyType       string `json:"bodyType"`
	FuelType       string `json:"fuelType"`
	EnginePowerHP  int    `json:"enginePowerHp"`
	EngineVolumeCC int    `json:"engineVolumeCc"`
	Transmission   string `json:"transmission"`
	DriveType      string `json:"driveType"`
}

type CarModelResponse struct {
	ID             string    `json:"id"`
	Brand          string    `json:"brand"`
	Model          string    `json:"model"`
	BasePrice      int64     `json:"basePrice"`
	BodyType       string    `json:"bodyType"`
	FuelType       string    `json:"fuelType"`
	EnginePowerHP  int       `json:"enginePowerHp"`
	EngineVolumeCC int       `json:"engineVolumeCc"`
	Transmission   string    `json:"transmission"`
	DriveType      string    `json:"driveType"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func CarModelToResponse(m car.CarModel) CarModelResponse {
	return CarModelResponse{
		ID:             m.ID,
		Brand:          m.Brand,
		Model:          m.Model,
		BasePrice:      m.BasePrice.Amount,
		BodyType:       string(m.BodyType),
		FuelType:       string(m.FuelType),
		EnginePowerHP:  m.EnginePowerHP,
		EngineVolumeCC: m.EngineVolumeCC,
		Transmission:   string(m.Transmission),
		DriveType:      string(m.DriveType),
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func CarModelsToResponse(models []car.CarModel) []CarModelResponse {
	result := make([]CarModelResponse, len(models))
	for i, m := range models {
		result[i] = CarModelToResponse(m)
	}
	return result
}
