package carserver

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/qxblvde/CarService/proto/storagepb"
	"github.com/qxblvde/CarService/storage-service/internal/domain"
)

type fakeInventory struct {
	cars    []domain.Car
	listErr error
}

func (f fakeInventory) ListAvailable(context.Context) ([]domain.Car, error) {
	return f.cars, f.listErr
}

func (f fakeInventory) GetByID(_ context.Context, id string) (domain.Car, error) {
	for _, c := range f.cars {
		if c.ID == id {
			return c, nil
		}
	}
	return domain.Car{}, domain.ErrCarNotFound
}

func TestListCars(t *testing.T) {
	srv := New(fakeInventory{cars: []domain.Car{
		{ID: "1", Brand: "BMW", Model: "X5", Year: 2023, Price: 5500000},
		{ID: "2", Brand: "Toyota", Model: "Camry", Year: 2022, Price: 2700000},
	}})

	resp, err := srv.ListCars(context.Background(), &storagepb.ListCarsRequest{})
	if err != nil {
		t.Fatalf("ListCars: %v", err)
	}
	if len(resp.GetCars()) != 2 {
		t.Fatalf("want 2 cars, got %d", len(resp.GetCars()))
	}
	if got := resp.GetCars()[0]; got.GetBrand() != "BMW" || got.GetPrice() != 5500000 {
		t.Fatalf("unexpected first car: %+v", got)
	}
}

func TestListCarsEmpty(t *testing.T) {
	srv := New(fakeInventory{})

	resp, err := srv.ListCars(context.Background(), &storagepb.ListCarsRequest{})
	if err != nil {
		t.Fatalf("ListCars: %v", err)
	}
	if len(resp.GetCars()) != 0 {
		t.Fatalf("want empty list, got %d", len(resp.GetCars()))
	}
}

func TestGetCarFound(t *testing.T) {
	srv := New(fakeInventory{cars: []domain.Car{
		{ID: "42", Brand: "Audi", Model: "A6", Year: 2021, Price: 3800000},
	}})

	car, err := srv.GetCar(context.Background(), &storagepb.GetCarRequest{Id: "42"})
	if err != nil {
		t.Fatalf("GetCar: %v", err)
	}
	if car.GetId() != "42" || car.GetModel() != "A6" {
		t.Fatalf("unexpected car: %+v", car)
	}
}

func TestGetCarNotFound(t *testing.T) {
	srv := New(fakeInventory{})

	_, err := srv.GetCar(context.Background(), &storagepb.GetCarRequest{Id: "missing"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound, got %v", err)
	}
}

func TestGetCarEmptyID(t *testing.T) {
	srv := New(fakeInventory{})

	_, err := srv.GetCar(context.Background(), &storagepb.GetCarRequest{})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("want InvalidArgument, got %v", err)
	}
}

func TestListCarsSourceError(t *testing.T) {
	srv := New(fakeInventory{listErr: errors.New("db down")})

	_, err := srv.ListCars(context.Background(), &storagepb.ListCarsRequest{})
	if status.Code(err) != codes.Internal {
		t.Fatalf("want Internal, got %v", err)
	}
}
