package order

import "testing"

func TestInStockTransitions(t *testing.T) {
	if err := ensureInStockTransition(InStockStatusPaid, InStockStatusReadyForPickUp); err != nil {
		t.Errorf("PAID -> READY_FOR_PICK_UP should be valid: %v", err)
	}
	if err := ensureInStockTransition(InStockStatusPaid, InStockStatusCanceled); err != nil {
		t.Errorf("PAID -> CANCELED should be valid: %v", err)
	}
	if err := ensureInStockTransition(InStockStatusReadyForPickUp, InStockStatusReadyForPickUp); err == nil {
		t.Error("READY_FOR_PICK_UP -> READY_FOR_PICK_UP should be invalid")
	}
}

func TestCustomApprovalPath(t *testing.T) {
	if err := ensureCustomTransition(CustomStatusPaid, CustomStatusWaitingDelivery); err != nil {
		t.Errorf("PAID -> WAITING_DELIVERY should be valid: %v", err)
	}
	if err := ensureCustomTransition(CustomStatusWaitingDelivery, CustomStatusReadyForPickUp); err != nil {
		t.Errorf("WAITING_DELIVERY -> READY_FOR_PICK_UP should be valid: %v", err)
	}
	if err := ensureCustomTransition(CustomStatusPaid, CustomStatusCanceled); err != nil {
		t.Errorf("PAID -> CANCELED should be valid: %v", err)
	}
}
