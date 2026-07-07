package service_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"

	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/application/service"
	"github.com/qxblvde/CarService/order-service/internal/domain/car"
	"github.com/qxblvde/CarService/order-service/internal/domain/order"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type fakeInStockRepo struct {
	orders map[string]order.InStockOrder
}

func (f *fakeInStockRepo) GetByID(_ context.Context, id string) (order.InStockOrder, error) {
	o, ok := f.orders[id]
	if !ok {
		return order.InStockOrder{}, apperrs.ErrOrderNotFound
	}
	return o, nil
}

func (f *fakeInStockRepo) List(context.Context, persistence.OrderFilter) ([]order.InStockOrder, error) {
	return nil, nil
}

func (f *fakeInStockRepo) Create(_ context.Context, o order.InStockOrder) error {
	f.orders[o.ID] = o
	return nil
}

func (f *fakeInStockRepo) UpdateStatus(_ context.Context, id string, status order.InStockStatus) error {
	o := f.orders[id]
	o.Status = status
	f.orders[id] = o
	return nil
}

func (f *fakeInStockRepo) UpdateStatusTx(ctx context.Context, _ pgx.Tx, id string, status order.InStockStatus) error {
	return f.UpdateStatus(ctx, id, status)
}

func paidInStockOrder(t *testing.T) order.InStockOrder {
	t.Helper()
	o, err := order.NewInStockOrder("order-1", "manager-1", "client-1",
		car.Car{ID: "car-1", Price: vo.MustMoney(1000)})
	if err != nil {
		t.Fatalf("new order: %v", err)
	}
	for _, step := range []func() error{o.Approve, o.WaitForPayment, o.MarkPaid} {
		if err := step(); err != nil {
			t.Fatalf("transition: %v", err)
		}
	}
	return o
}

func newService(repo *fakeInStockRepo) *service.OrderService {
	return service.NewOrderService(nil, repo, nil, nil, nil, nil, nil, nil)
}

func TestApplyApprovalReadyForPickup(t *testing.T) {
	o := paidInStockOrder(t)
	repo := &fakeInStockRepo{orders: map[string]order.InStockOrder{o.ID: o}}
	svc := newService(repo)

	if err := svc.ApplyApproval(context.Background(), service.OrderTypeInStock, o.ID); err != nil {
		t.Fatalf("ApplyApproval: %v", err)
	}
	if got := repo.orders[o.ID].Status; got != order.InStockStatusReadyForPickUp {
		t.Fatalf("want READY_FOR_PICK_UP, got %s", got)
	}
}

func TestApplyRejectionCancels(t *testing.T) {
	o := paidInStockOrder(t)
	repo := &fakeInStockRepo{orders: map[string]order.InStockOrder{o.ID: o}}
	svc := newService(repo)

	if err := svc.ApplyRejection(context.Background(), service.OrderTypeInStock, o.ID); err != nil {
		t.Fatalf("ApplyRejection: %v", err)
	}
	if got := repo.orders[o.ID].Status; got != order.InStockStatusCanceled {
		t.Fatalf("want CANCELED, got %s", got)
	}
}

func TestApplyApprovalIsNotRepeatable(t *testing.T) {
	o := paidInStockOrder(t)
	repo := &fakeInStockRepo{orders: map[string]order.InStockOrder{o.ID: o}}
	svc := newService(repo)

	if err := svc.ApplyApproval(context.Background(), service.OrderTypeInStock, o.ID); err != nil {
		t.Fatalf("first ApplyApproval: %v", err)
	}
	if err := svc.ApplyApproval(context.Background(), service.OrderTypeInStock, o.ID); err == nil {
		t.Fatal("second ApplyApproval should fail (already ready)")
	}
}
