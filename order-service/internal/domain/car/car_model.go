package car

import (
	"time"

	"github.com/qxblvde/CarService/internal/domain/errs"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type CarModel struct {
	ID             string
	Brand          string
	Model          string
	BasePrice      vo.Money
	BodyType       BodyType
	FuelType       FuelType
	EnginePowerHP  int
	EngineVolumeCC int
	Transmission   TransmissionType
	DriveType      DriveType
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

func NewCarModel(
	id, brand, model string,
	basePrice vo.Money,
	bodyType BodyType,
	fuelType FuelType,
	enginePowerHP int,
	engineVolumeCC int,
	transmission TransmissionType,
	driveType DriveType,
) (CarModel, error) {
	if id == "" {
		return CarModel{}, errs.Validation("car model id is required")
	}
	if brand == "" {
		return CarModel{}, errs.Validation("car brand is required")
	}
	if model == "" {
		return CarModel{}, errs.Validation("car model name is required")
	}
	if basePrice.IsNegative() {
		return CarModel{}, errs.Validation("car model base price must be non-negative")
	}
	if !bodyType.IsValid() {
		return CarModel{}, errs.Validation("invalid body type")
	}
	if !fuelType.IsValid() {
		return CarModel{}, errs.Validation("invalid fuel type")
	}
	if !transmission.IsValid() {
		return CarModel{}, errs.Validation("invalid transmission type")
	}
	if !driveType.IsValid() {
		return CarModel{}, errs.Validation("invalid drive type")
	}
	if enginePowerHP <= 0 {
		return CarModel{}, errs.Validation("engine power must be positive")
	}
	if engineVolumeCC <= 0 {
		return CarModel{}, errs.Validation("engine volume must be positive")
	}

	now := time.Now().UTC()
	return CarModel{
		ID:             id,
		Brand:          brand,
		Model:          model,
		BasePrice:      basePrice,
		BodyType:       bodyType,
		FuelType:       fuelType,
		EnginePowerHP:  enginePowerHP,
		EngineVolumeCC: engineVolumeCC,
		Transmission:   transmission,
		DriveType:      driveType,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (m CarModel) DisplayName() string {
	if m.Brand == "" {
		return m.Model
	}
	if m.Model == "" {
		return m.Brand
	}
	return m.Brand + " " + m.Model
}
