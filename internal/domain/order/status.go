package order

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
