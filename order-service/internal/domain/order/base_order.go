package order

import (
	"time"

	"github.com/qxblvde/CarService/internal/domain/errs"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type OrderBase struct {
	ID        string
	ManagerID string
	ClientID  string
	Price     vo.Money
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func newOrderBase(id, managerID, clientID string, price vo.Money) (OrderBase, error) {
	if id == "" {
		return OrderBase{}, errs.Validation("order id is required")
	}
	if managerID == "" {
		return OrderBase{}, errs.Validation("manager id is required")
	}
	if clientID == "" {
		return OrderBase{}, errs.Validation("client id is required")
	}
	if price.IsNegative() {
		return OrderBase{}, errs.Validation("order price must be non-negative")
	}

	now := time.Now().UTC()
	return OrderBase{
		ID:        id,
		ManagerID: managerID,
		ClientID:  clientID,
		Price:     price,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
