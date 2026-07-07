package order

import (
	"time"

	"github.com/qxblvde/CarService/internal/domain/errs"
	"github.com/qxblvde/CarService/order-service/internal/domain/car"
)

type CustomOrder struct {
	OrderBase
	CarModelID    string
	Configuration car.Configuration
	Status        CustomStatus
}

func NewCustomOrder(id, managerID, clientID string, model car.CarModel, configuration car.Configuration) (CustomOrder, error) {
	if model.ID == "" {
		return CustomOrder{}, errs.Validation("car model is required")
	}
	if configuration.CarModelID != model.ID {
		return CustomOrder{}, errs.Validation("configuration belongs to a different car model")
	}
	if len(configuration.Parts) == 0 {
		return CustomOrder{}, errs.Validation("configuration must have at least one part")
	}

	price := configuration.TotalPrice(model.BasePrice)

	base, err := newOrderBase(id, managerID, clientID, price)
	if err != nil {
		return CustomOrder{}, err
	}

	return CustomOrder{
		OrderBase:     base,
		CarModelID:    model.ID,
		Configuration: configuration,
		Status:        CustomStatusCreated,
	}, nil
}

func (o *CustomOrder) SetStatus(status CustomStatus) error {
	if err := ensureCustomTransition(o.Status, status); err != nil {
		return err
	}
	o.Status = status
	o.UpdatedAt = time.Now().UTC()
	return nil
}

func (o *CustomOrder) Approve() error {
	return o.SetStatus(CustomStatusApproved)
}

func (o *CustomOrder) WaitForPayment() error {
	return o.SetStatus(CustomStatusWaitingPayment)
}

func (o *CustomOrder) MarkPaid() error {
	return o.SetStatus(CustomStatusPaid)
}

func (o *CustomOrder) WaitingDelivery() error {
	return o.SetStatus(CustomStatusWaitingDelivery)
}

func (o *CustomOrder) ReadyForPickup() error {
	return o.SetStatus(CustomStatusReadyForPickUp)
}

func (o *CustomOrder) Complete() error {
	return o.SetStatus(CustomStatusCompleted)
}

func (o *CustomOrder) Cancel() error {
	return o.SetStatus(CustomStatusCanceled)
}
