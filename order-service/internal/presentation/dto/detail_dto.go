package dto

import (
	"time"

	"github.com/qxblvde/CarService/order-service/internal/domain/detail"
)

type CreateDetailRequest struct {
	Name                  string   `json:"name"`
	Type                  string   `json:"type"`
	Price                 int64    `json:"price"`
	CompatibleCarModelIDs []string `json:"compatibleCarModelIds"`
}

type DetailResponse struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Type                  string    `json:"type"`
	Price                 int64     `json:"price"`
	CompatibleCarModelIDs []string  `json:"compatibleCarModelIds"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

func DetailToResponse(d detail.Detail) DetailResponse {
	ids := d.CompatibleCarModelIDs
	if ids == nil {
		ids = []string{}
	}
	return DetailResponse{
		ID:                    d.ID,
		Name:                  d.Name,
		Type:                  string(d.Type),
		Price:                 d.Price.Amount,
		CompatibleCarModelIDs: ids,
		CreatedAt:             d.CreatedAt,
		UpdatedAt:             d.UpdatedAt,
	}
}

func DetailsToResponse(details []detail.Detail) []DetailResponse {
	result := make([]DetailResponse, len(details))
	for i, d := range details {
		result[i] = DetailToResponse(d)
	}
	return result
}
