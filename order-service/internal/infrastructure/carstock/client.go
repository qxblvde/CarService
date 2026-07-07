package carstock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/qxblvde/CarService/internal/grpclog"
	"github.com/qxblvde/CarService/proto/storagepb"
)

var ErrUnavailable = errors.New("storage service unavailable")

var ErrCarNotFound = errors.New("car not found")

type Car struct {
	ID    string
	Brand string
	Model string
	Year  int
	Price int64
}

type Client struct {
	conn    *grpc.ClientConn
	cars    storagepb.CarStockClient
	timeout time.Duration
}

func Dial(addr string, timeout time.Duration, opts ...grpc.DialOption) (*Client, error) {
	dialOpts := append([]grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpclog.UnaryClient),
	}, opts...)

	conn, err := grpc.NewClient(addr, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("dial storage: %w", err)
	}
	return &Client{
		conn:    conn,
		cars:    storagepb.NewCarStockClient(conn),
		timeout: timeout,
	}, nil
}

func (c *Client) List(ctx context.Context) ([]Car, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.cars.ListCars(ctx, &storagepb.ListCarsRequest{})
	if err != nil {
		return nil, mapError(err)
	}

	cars := make([]Car, 0, len(resp.GetCars()))
	for _, c := range resp.GetCars() {
		cars = append(cars, fromProto(c))
	}
	return cars, nil
}

func (c *Client) Get(ctx context.Context, id string) (Car, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	car, err := c.cars.GetCar(ctx, &storagepb.GetCarRequest{Id: id})
	if err != nil {
		return Car{}, mapError(err)
	}
	return fromProto(car), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func fromProto(c *storagepb.Car) Car {
	return Car{
		ID:    c.GetId(),
		Brand: c.GetBrand(),
		Model: c.GetModel(),
		Year:  int(c.GetYear()),
		Price: c.GetPrice(),
	}
}

func mapError(err error) error {
	switch status.Code(err) {
	case codes.NotFound:
		return ErrCarNotFound
	case codes.Unavailable, codes.DeadlineExceeded:
		return fmt.Errorf("%w: %v", ErrUnavailable, err)
	default:
		return err
	}
}
