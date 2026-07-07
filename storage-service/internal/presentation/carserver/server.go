package carserver

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/qxblvde/CarService/proto/storagepb"
	"github.com/qxblvde/CarService/storage-service/internal/domain"
)

type Inventory interface {
	ListAvailable(ctx context.Context) ([]domain.Car, error)
	GetByID(ctx context.Context, id string) (domain.Car, error)
}

type Server struct {
	storagepb.UnimplementedCarStockServer
	cars Inventory
}

func New(cars Inventory) *Server {
	return &Server{cars: cars}
}

func (s *Server) ListCars(ctx context.Context, _ *storagepb.ListCarsRequest) (*storagepb.ListCarsResponse, error) {
	cars, err := s.cars.ListAvailable(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &storagepb.ListCarsResponse{Cars: make([]*storagepb.Car, 0, len(cars))}
	for _, c := range cars {
		resp.Cars = append(resp.Cars, toProto(c))
	}
	return resp, nil
}

func (s *Server) GetCar(ctx context.Context, req *storagepb.GetCarRequest) (*storagepb.Car, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	c, err := s.cars.GetByID(ctx, req.GetId())
	if errors.Is(err, domain.ErrCarNotFound) {
		return nil, status.Errorf(codes.NotFound, "car %s not found", req.GetId())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProto(c), nil
}

func toProto(c domain.Car) *storagepb.Car {
	return &storagepb.Car{
		Id:    c.ID,
		Brand: c.Brand,
		Model: c.Model,
		Year:  int32(c.Year),
		Price: c.Price,
	}
}
