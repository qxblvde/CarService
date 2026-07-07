package carstock_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/qxblvde/CarService/order-service/internal/infrastructure/carstock"
	"github.com/qxblvde/CarService/proto/storagepb"
)

type stubServer struct {
	storagepb.UnimplementedCarStockServer
	cars  []*storagepb.Car
	delay time.Duration
}

func (s *stubServer) ListCars(ctx context.Context, _ *storagepb.ListCarsRequest) (*storagepb.ListCarsResponse, error) {
	if s.delay > 0 {
		select {
		case <-time.After(s.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return &storagepb.ListCarsResponse{Cars: s.cars}, nil
}

func (s *stubServer) GetCar(_ context.Context, req *storagepb.GetCarRequest) (*storagepb.Car, error) {
	for _, c := range s.cars {
		if c.GetId() == req.GetId() {
			return c, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "car %s not found", req.GetId())
}

func startServer(t *testing.T, stub *stubServer, timeout time.Duration) (*carstock.Client, *grpc.Server) {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	gs := grpc.NewServer()
	storagepb.RegisterCarStockServer(gs, stub)
	go func() { _ = gs.Serve(lis) }()

	dialer := grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.DialContext(ctx)
	})
	client, err := carstock.Dial("passthrough:///bufnet", timeout,
		dialer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	t.Cleanup(func() {
		_ = client.Close()
		gs.Stop()
	})
	return client, gs
}

func sampleCars() []*storagepb.Car {
	return []*storagepb.Car{
		{Id: "1", Brand: "BMW", Model: "X5", Year: 2023, Price: 5500000},
		{Id: "2", Brand: "Toyota", Model: "Camry", Year: 2022, Price: 2700000},
	}
}

func TestClientList(t *testing.T) {
	client, _ := startServer(t, &stubServer{cars: sampleCars()}, time.Second)

	cars, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(cars) != 2 {
		t.Fatalf("want 2 cars, got %d", len(cars))
	}
	if cars[0].Brand != "BMW" || cars[0].Price != 5500000 {
		t.Fatalf("unexpected car: %+v", cars[0])
	}
}

func TestClientListEmpty(t *testing.T) {
	client, _ := startServer(t, &stubServer{}, time.Second)

	cars, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(cars) != 0 {
		t.Fatalf("want empty result, got %d", len(cars))
	}
}

func TestClientGet(t *testing.T) {
	client, _ := startServer(t, &stubServer{cars: sampleCars()}, time.Second)

	car, err := client.Get(context.Background(), "2")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if car.Model != "Camry" {
		t.Fatalf("unexpected car: %+v", car)
	}
}

func TestClientGetNotFound(t *testing.T) {
	client, _ := startServer(t, &stubServer{cars: sampleCars()}, time.Second)

	_, err := client.Get(context.Background(), "does-not-exist")
	if !errors.Is(err, carstock.ErrCarNotFound) {
		t.Fatalf("want ErrCarNotFound, got %v", err)
	}
}

func TestClientTimeout(t *testing.T) {
	client, _ := startServer(t, &stubServer{cars: sampleCars(), delay: 300 * time.Millisecond}, 50*time.Millisecond)

	_, err := client.List(context.Background())
	if !errors.Is(err, carstock.ErrUnavailable) {
		t.Fatalf("want ErrUnavailable on timeout, got %v", err)
	}
}

func TestClientUnavailable(t *testing.T) {
	client, gs := startServer(t, &stubServer{cars: sampleCars()}, time.Second)

	gs.Stop()

	_, err := client.List(context.Background())
	if !errors.Is(err, carstock.ErrUnavailable) {
		t.Fatalf("want ErrUnavailable when server is down, got %v", err)
	}
}
