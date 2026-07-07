package order

import (
	"time"

	"github.com/qxblvde/CarService/internal/domain/errs"
	"github.com/qxblvde/CarService/order-service/internal/domain/car"
)

type InStockOrder struct {
	OrderBase
	CarID  string
	Status InStockStatus
}

func NewInStockOrder(id, managerID, clientID string, carEntity car.Car) (InStockOrder, error) {
	if carEntity.ID == "" {
		return InStockOrder{}, errs.Validation("car is required")
	}

	base, err := newOrderBase(id, managerID, clientID, carEntity.Price)
	if err != nil {
		return InStockOrder{}, err
	}

	return InStockOrder{
		OrderBase: base,
		CarID:     carEntity.ID,
		Status:    InStockStatusCreated,
	}, nil
}

func (o *InStockOrder) SetStatus(status InStockStatus) error {
	if err := ensureInStockTransition(o.Status, status); err != nil {
		return err
	}
	o.Status = status
	o.UpdatedAt = time.Now().UTC()
	return nil
}

func (o *InStockOrder) Approve() error {
	return o.SetStatus(InStockStatusApproved)
}

func (o *InStockOrder) WaitForPayment() error {
	return o.SetStatus(InStockStatusWaitingPayment)
}

func (o *InStockOrder) MarkPaid() error {
	return o.SetStatus(InStockStatusPaid)
}

func (o *InStockOrder) ReadyForPickup() error {
	return o.SetStatus(InStockStatusReadyForPickUp)
}

func (o *InStockOrder) Complete() error {
	return o.SetStatus(InStockStatusCompleted)
}

func (o *InStockOrder) Cancel() error {
	return o.SetStatus(InStockStatusCanceled)
}
