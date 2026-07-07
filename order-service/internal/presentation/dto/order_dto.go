package dto

import (
	"time"

	"github.com/qxblvde/CarService/order-service/internal/domain/order"
)

type CreateInStockOrderRequest struct {
	ManagerID string `json:"managerId"`
	ClientID  string `json:"clientId"`
	CarID     string `json:"carId"`
}

type InStockOrderResponse struct {
	ID        string    `json:"id"`
	ManagerID string    `json:"managerId"`
	ClientID  string    `json:"clientId"`
	CarID     string    `json:"carId"`
	Status    string    `json:"status"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func InStockOrderToResponse(o order.InStockOrder) InStockOrderResponse {
	return InStockOrderResponse{
		ID:        o.ID,
		ManagerID: o.ManagerID,
		ClientID:  o.ClientID,
		CarID:     o.CarID,
		Status:    string(o.Status),
		Price:     o.Price.Amount,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func InStockOrdersToResponse(orders []order.InStockOrder) []InStockOrderResponse {
	result := make([]InStockOrderResponse, len(orders))
	for i, o := range orders {
		result[i] = InStockOrderToResponse(o)
	}
	return result
}

type CreateCustomOrderRequest struct {
	ManagerID  string   `json:"managerId"`
	ClientID   string   `json:"clientId"`
	CarModelID string   `json:"carModelId"`
	DetailIDs  []string `json:"detailIds"`
}

type ConfigurationPartResponse struct {
	DetailID string `json:"detailId"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Price    int64  `json:"price"`
}

type CustomOrderResponse struct {
	ID         string                      `json:"id"`
	ManagerID  string                      `json:"managerId"`
	ClientID   string                      `json:"clientId"`
	CarModelID string                      `json:"carModelId"`
	Status     string                      `json:"status"`
	Price      int64                       `json:"price"`
	Parts      []ConfigurationPartResponse `json:"parts"`
	CreatedAt  time.Time                   `json:"createdAt"`
	UpdatedAt  time.Time                   `json:"updatedAt"`
}

func CustomOrderToResponse(o order.CustomOrder) CustomOrderResponse {
	parts := make([]ConfigurationPartResponse, 0, len(o.Configuration.Parts))
	for _, p := range o.Configuration.Parts {
		parts = append(parts, ConfigurationPartResponse{
			DetailID: p.ID,
			Name:     p.Name,
			Type:     string(p.Type),
			Price:    p.Price.Amount,
		})
	}
	return CustomOrderResponse{
		ID:         o.ID,
		ManagerID:  o.ManagerID,
		ClientID:   o.ClientID,
		CarModelID: o.CarModelID,
		Status:     string(o.Status),
		Price:      o.Price.Amount,
		Parts:      parts,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}

func CustomOrdersToResponse(orders []order.CustomOrder) []CustomOrderResponse {
	result := make([]CustomOrderResponse, len(orders))
	for i, o := range orders {
		result[i] = CustomOrderToResponse(o)
	}
	return result
}
