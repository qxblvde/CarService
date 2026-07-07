package order

import "github.com/qxblvde/CarService/internal/domain/errs"

type CustomStatus string
type InStockStatus string

const (
	CustomStatusCreated         CustomStatus = "CREATED"
	CustomStatusApproved        CustomStatus = "APPROVED"
	CustomStatusWaitingPayment  CustomStatus = "WAITING_PAYMENT"
	CustomStatusPaid            CustomStatus = "PAID"
	CustomStatusWaitingDelivery CustomStatus = "WAITING_DELIVERY"
	CustomStatusReadyForPickUp  CustomStatus = "READY_FOR_PICK_UP"
	CustomStatusCompleted       CustomStatus = "COMPLETED"
	CustomStatusCanceled        CustomStatus = "CANCELED"
)

const (
	InStockStatusCreated        InStockStatus = "CREATED"
	InStockStatusApproved       InStockStatus = "APPROVED"
	InStockStatusWaitingPayment InStockStatus = "WAITING_PAYMENT"
	InStockStatusPaid           InStockStatus = "PAID"
	InStockStatusReadyForPickUp InStockStatus = "READY_FOR_PICK_UP"
	InStockStatusCompleted      InStockStatus = "COMPLETED"
	InStockStatusCanceled       InStockStatus = "CANCELED"
)

func ensureCustomTransition(current, next CustomStatus) error {
	switch current {
	case CustomStatusCreated:
		if next == CustomStatusApproved || next == CustomStatusCanceled {
			return nil
		}
	case CustomStatusApproved:
		if next == CustomStatusWaitingPayment || next == CustomStatusCanceled {
			return nil
		}
	case CustomStatusWaitingPayment:
		if next == CustomStatusPaid || next == CustomStatusCanceled {
			return nil
		}
	case CustomStatusPaid:
		if next == CustomStatusWaitingDelivery || next == CustomStatusCanceled {
			return nil
		}
	case CustomStatusWaitingDelivery:
		if next == CustomStatusReadyForPickUp || next == CustomStatusCanceled {
			return nil
		}
	case CustomStatusReadyForPickUp:
		if next == CustomStatusCompleted || next == CustomStatusCanceled {
			return nil
		}
	case CustomStatusCompleted, CustomStatusCanceled:
		return errs.InvalidTransition
	}
	return errs.Validation("invalid custom order status")
}

func ensureInStockTransition(current, next InStockStatus) error {
	switch current {
	case InStockStatusCreated:
		if next == InStockStatusApproved || next == InStockStatusCanceled {
			return nil
		}
	case InStockStatusApproved:
		if next == InStockStatusWaitingPayment || next == InStockStatusCanceled {
			return nil
		}
	case InStockStatusWaitingPayment:
		if next == InStockStatusPaid || next == InStockStatusCanceled {
			return nil
		}
	case InStockStatusPaid:
		if next == InStockStatusReadyForPickUp || next == InStockStatusCanceled {
			return nil
		}
	case InStockStatusReadyForPickUp:
		if next == InStockStatusCompleted || next == InStockStatusCanceled {
			return nil
		}
	case InStockStatusCompleted, InStockStatusCanceled:
		return errs.InvalidTransition
	}
	return errs.Validation("invalid in-stock order status")
}
