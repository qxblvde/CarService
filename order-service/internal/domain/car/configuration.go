package car

import (
	"github.com/qxblvde/CarService/internal/domain/errs"
	"github.com/qxblvde/CarService/order-service/internal/domain/detail"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type Configuration struct {
	CarModelID string
	Parts      map[detail.DetailType]detail.Detail
}

func NewConfiguration(carModelID string, parts ...detail.Detail) (Configuration, error) {
	if carModelID == "" {
		return Configuration{}, errs.Validation("car model id is required")
	}

	cfg := Configuration{
		CarModelID: carModelID,
		Parts:      make(map[detail.DetailType]detail.Detail, len(parts)),
	}

	for _, part := range parts {
		if !part.Type.IsValid() {
			return Configuration{}, errs.Validation("invalid part type: " + part.Type.String())
		}
		if _, exists := cfg.Parts[part.Type]; exists {
			return Configuration{}, errs.DuplicateSelection(part.Type.String())
		}
		cfg.Parts[part.Type] = part
	}

	return cfg, nil
}

func (c Configuration) TotalPrice(basePrice vo.Money) vo.Money {
	total := basePrice
	for _, part := range c.Parts {
		total = total.Add(part.Price.Amount)
	}
	return total
}
