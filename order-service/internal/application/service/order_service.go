package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/domain/car"
	"github.com/qxblvde/CarService/order-service/internal/domain/detail"
	"github.com/qxblvde/CarService/order-service/internal/domain/order"
	"github.com/qxblvde/CarService/order-service/internal/infrastructure"
)

const (
	OrderTypeInStock = "IN_STOCK"
	OrderTypeCustom  = "CUSTOM"
)

type OrderSentForApproval struct {
	OrderID    string    `json:"orderId"`
	OrderType  string    `json:"orderType"`
	TraceID    string    `json:"traceId"`
	OccurredAt time.Time `json:"occurredAt"`
}

type OrderService struct {
	customOrders  persistence.CustomOrderRepository
	inStockOrders persistence.InStockOrderRepository
	cars          persistence.CarRepository
	models        persistence.CarModelRepository
	details       persistence.DetailRepository
	users         persistence.UserRepository
	outbox        *infrastructure.OutboxRepository
	pool          *pgxpool.Pool
}

func NewOrderService(
	customOrders persistence.CustomOrderRepository,
	inStockOrders persistence.InStockOrderRepository,
	cars persistence.CarRepository,
	models persistence.CarModelRepository,
	details persistence.DetailRepository,
	users persistence.UserRepository,
	outbox *infrastructure.OutboxRepository,
	pool *pgxpool.Pool,
) *OrderService {
	return &OrderService{
		customOrders:  customOrders,
		inStockOrders: inStockOrders,
		cars:          cars,
		models:        models,
		details:       details,
		users:         users,
		outbox:        outbox,
		pool:          pool,
	}
}

type CreateInStockOrderParams struct {
	ManagerID string
	ClientID  string
	CarID     string
}

func (s *OrderService) CreateInStockOrder(ctx context.Context, p CreateInStockOrderParams) (order.InStockOrder, error) {
	if _, err := s.users.GetByID(ctx, p.ManagerID); err != nil {
		return order.InStockOrder{}, apperrs.ErrUserNotFound
	}
	if _, err := s.users.GetByID(ctx, p.ClientID); err != nil {
		return order.InStockOrder{}, apperrs.ErrUserNotFound
	}

	carEntity, err := s.cars.GetByID(ctx, p.CarID)
	if err != nil {
		return order.InStockOrder{}, apperrs.ErrCarNotFound
	}

	o, err := order.NewInStockOrder(uuid.NewString(), p.ManagerID, p.ClientID, carEntity)
	if err != nil {
		return order.InStockOrder{}, err
	}

	if err := s.inStockOrders.Create(ctx, o); err != nil {
		return order.InStockOrder{}, err
	}
	return o, nil
}

func (s *OrderService) GetInStockOrder(ctx context.Context, id string) (order.InStockOrder, error) {
	return s.inStockOrders.GetByID(ctx, id)
}

func (s *OrderService) ListInStockOrders(ctx context.Context, filter persistence.OrderFilter) ([]order.InStockOrder, error) {
	return s.inStockOrders.List(ctx, filter)
}

func (s *OrderService) transitionInStock(ctx context.Context, id string, apply func(*order.InStockOrder) error) error {
	o, err := s.inStockOrders.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := apply(&o); err != nil {
		return err
	}
	return s.inStockOrders.UpdateStatus(ctx, id, o.Status)
}

func (s *OrderService) ApproveInStockOrder(ctx context.Context, id string) error {
	return s.transitionInStock(ctx, id, (*order.InStockOrder).Approve)
}

func (s *OrderService) InStockOrderWaitPayment(ctx context.Context, id string) error {
	return s.transitionInStock(ctx, id, (*order.InStockOrder).WaitForPayment)
}

func (s *OrderService) MarkInStockOrderPaid(ctx context.Context, id string) error {
	o, err := s.inStockOrders.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := o.MarkPaid(); err != nil {
		return err
	}
	return infrastructure.WithTx(ctx, s.pool, func(ctx context.Context, tx infrastructure.Tx) error {
		if err := s.inStockOrders.UpdateStatusTx(ctx, tx, id, o.Status); err != nil {
			return err
		}
		event := OrderSentForApproval{
			OrderID:    id,
			OrderType:  "IN_STOCK",
			TraceID:    uuid.NewString(),
			OccurredAt: time.Now().UTC(),
		}
		return s.outbox.SaveTx(ctx, tx, "order.sent_for_approval", event)
	})
}

func (s *OrderService) InStockOrderReadyForPickup(ctx context.Context, id string) error {
	return s.transitionInStock(ctx, id, (*order.InStockOrder).ReadyForPickup)
}

func (s *OrderService) CompleteInStockOrder(ctx context.Context, id string) error {
	return s.transitionInStock(ctx, id, (*order.InStockOrder).Complete)
}

func (s *OrderService) CancelInStockOrder(ctx context.Context, id string) error {
	return s.transitionInStock(ctx, id, (*order.InStockOrder).Cancel)
}

type CreateCustomOrderParams struct {
	ManagerID  string
	ClientID   string
	CarModelID string
	DetailIDs  []string
}

func (s *OrderService) CreateCustomOrder(ctx context.Context, p CreateCustomOrderParams) (order.CustomOrder, error) {
	if _, err := s.users.GetByID(ctx, p.ManagerID); err != nil {
		return order.CustomOrder{}, apperrs.ErrUserNotFound
	}
	if _, err := s.users.GetByID(ctx, p.ClientID); err != nil {
		return order.CustomOrder{}, apperrs.ErrUserNotFound
	}

	model, err := s.models.GetByID(ctx, p.CarModelID)
	if err != nil {
		return order.CustomOrder{}, apperrs.ErrNotFound
	}

	parts, err := s.loadAndValidateParts(ctx, p.CarModelID, p.DetailIDs)
	if err != nil {
		return order.CustomOrder{}, err
	}

	cfg, err := car.NewConfiguration(p.CarModelID, parts...)
	if err != nil {
		return order.CustomOrder{}, err
	}

	o, err := order.NewCustomOrder(uuid.NewString(), p.ManagerID, p.ClientID, model, cfg)
	if err != nil {
		return order.CustomOrder{}, err
	}

	if err := s.customOrders.Create(ctx, o); err != nil {
		return order.CustomOrder{}, err
	}
	return o, nil
}

func (s *OrderService) loadAndValidateParts(ctx context.Context, carModelID string, detailIDs []string) ([]detail.Detail, error) {
	parts := make([]detail.Detail, 0, len(detailIDs))
	for _, id := range detailIDs {
		d, err := s.details.GetByID(ctx, id)
		if err != nil {
			return nil, apperrs.ErrNotFound
		}
		if !d.CompatibleWith(carModelID) {
			return nil, apperrs.ErrPartNotCompatible
		}
		parts = append(parts, d)
	}
	return parts, nil
}

func (s *OrderService) GetCustomOrder(ctx context.Context, id string) (order.CustomOrder, error) {
	return s.customOrders.GetByID(ctx, id)
}

func (s *OrderService) ListCustomOrders(ctx context.Context, filter persistence.OrderFilter) ([]order.CustomOrder, error) {
	return s.customOrders.List(ctx, filter)
}

func (s *OrderService) transitionCustom(ctx context.Context, id string, apply func(*order.CustomOrder) error) error {
	o, err := s.customOrders.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := apply(&o); err != nil {
		return err
	}
	return s.customOrders.UpdateStatus(ctx, id, o.Status)
}

func (s *OrderService) ApproveCustomOrder(ctx context.Context, id string) error {
	return s.transitionCustom(ctx, id, (*order.CustomOrder).Approve)
}

func (s *OrderService) CustomOrderWaitPayment(ctx context.Context, id string) error {
	return s.transitionCustom(ctx, id, (*order.CustomOrder).WaitForPayment)
}

func (s *OrderService) MarkCustomOrderPaid(ctx context.Context, id string) error {
	o, err := s.customOrders.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := o.MarkPaid(); err != nil {
		return err
	}
	return infrastructure.WithTx(ctx, s.pool, func(ctx context.Context, tx infrastructure.Tx) error {
		if err := s.customOrders.UpdateStatusTx(ctx, tx, id, o.Status); err != nil {
			return err
		}
		event := OrderSentForApproval{
			OrderID:    id,
			OrderType:  "CUSTOM",
			TraceID:    uuid.NewString(),
			OccurredAt: time.Now().UTC(),
		}
		return s.outbox.SaveTx(ctx, tx, "order.sent_for_approval", event)
	})
}

func (s *OrderService) CustomOrderWaitingDelivery(ctx context.Context, id string) error {
	return s.transitionCustom(ctx, id, (*order.CustomOrder).WaitingDelivery)
}

func (s *OrderService) CustomOrderReadyForPickup(ctx context.Context, id string) error {
	return s.transitionCustom(ctx, id, (*order.CustomOrder).ReadyForPickup)
}

func (s *OrderService) CompleteCustomOrder(ctx context.Context, id string) error {
	return s.transitionCustom(ctx, id, (*order.CustomOrder).Complete)
}

func (s *OrderService) CancelCustomOrder(ctx context.Context, id string) error {
	return s.transitionCustom(ctx, id, (*order.CustomOrder).Cancel)
}

func (s *OrderService) ApplyApproval(ctx context.Context, orderType, id string) error {
	switch orderType {
	case OrderTypeInStock:
		return s.InStockOrderReadyForPickup(ctx, id)
	case OrderTypeCustom:
		if err := s.CustomOrderWaitingDelivery(ctx, id); err != nil {
			return err
		}
		return s.CustomOrderReadyForPickup(ctx, id)
	default:
		return fmt.Errorf("unknown order type %q", orderType)
	}
}

func (s *OrderService) ApplyRejection(ctx context.Context, orderType, id string) error {
	switch orderType {
	case OrderTypeInStock:
		return s.CancelInStockOrder(ctx, id)
	case OrderTypeCustom:
		return s.CancelCustomOrder(ctx, id)
	default:
		return fmt.Errorf("unknown order type %q", orderType)
	}
}
