package persistence

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/qxblvde/CarService/order-service/internal/domain/order"
	"github.com/qxblvde/CarService/order-service/internal/domain/test_drive"
)

type OrderFilter struct {
	ClientID  string
	ManagerID string
}

type TestDriveFilter struct {
	ClientID string
	CarID    string
}

type CustomOrderRepository interface {
	GetByID(ctx context.Context, id string) (order.CustomOrder, error)
	List(ctx context.Context, filter OrderFilter) ([]order.CustomOrder, error)
	Create(ctx context.Context, o order.CustomOrder) error
	UpdateStatus(ctx context.Context, id string, status order.CustomStatus) error
	UpdateStatusTx(ctx context.Context, tx pgx.Tx, id string, status order.CustomStatus) error
}

type InStockOrderRepository interface {
	GetByID(ctx context.Context, id string) (order.InStockOrder, error)
	List(ctx context.Context, filter OrderFilter) ([]order.InStockOrder, error)
	Create(ctx context.Context, o order.InStockOrder) error
	UpdateStatus(ctx context.Context, id string, status order.InStockStatus) error
	UpdateStatusTx(ctx context.Context, tx pgx.Tx, id string, status order.InStockStatus) error
}

type TestDriveRepository interface {
	GetByID(ctx context.Context, id string) (test_drive.TestDriveRequest, error)
	List(ctx context.Context, filter TestDriveFilter) ([]test_drive.TestDriveRequest, error)
	Create(ctx context.Context, r test_drive.TestDriveRequest) error
	Cancel(ctx context.Context, id string) error
}
