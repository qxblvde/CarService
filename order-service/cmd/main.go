package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/qxblvde/CarService/internal/broker"
	"github.com/qxblvde/CarService/internal/config"
	"github.com/qxblvde/CarService/internal/middleware"
	"github.com/qxblvde/CarService/internal/server"
	"github.com/qxblvde/CarService/order-service/internal/application/service"
	"github.com/qxblvde/CarService/order-service/internal/infrastructure"
	"github.com/qxblvde/CarService/order-service/internal/infrastructure/carstock"
	"github.com/qxblvde/CarService/order-service/internal/infrastructure/messaging"
	"github.com/qxblvde/CarService/order-service/internal/infrastructure/migrations"
	"github.com/qxblvde/CarService/order-service/internal/presentation/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := migrations.Run(cfg.DBURL); err != nil {
		log.Fatalf("migrations error: %v", err)
	}
	log.Println("migrations applied")

	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("db ping error: %v", err)
	}
	log.Println("connected to database")

	storageClient, err := carstock.Dial(cfg.StorageGRPCAddr, 3*time.Second)
	if err != nil {
		log.Fatalf("storage client error: %v", err)
	}
	defer storageClient.Close()
	log.Printf("storage gRPC client dialing %s", cfg.StorageGRPCAddr)

	outboxRepo := infrastructure.NewOutboxRepository(pool)

	carRepo := infrastructure.NewCarRepository(pool)
	carModelRepo := infrastructure.NewCarModelRepository(pool)
	detailRepo := infrastructure.NewDetailRepository(pool)
	customOrderRepo := infrastructure.NewCustomOrderRepository(pool)
	inStockOrderRepo := infrastructure.NewInStockOrderRepository(pool)
	testDriveRepo := infrastructure.NewTestDriveRepository(pool)
	userRepo := infrastructure.NewUserRepository(pool)

	carModelSvc := service.NewCarModelService(carModelRepo)
	detailSvc := service.NewDetailService(detailRepo)
	orderSvc := service.NewOrderService(customOrderRepo, inStockOrderRepo, carRepo, carModelRepo, detailRepo, userRepo, outboxRepo, pool)
	testDriveSvc := service.NewTestDriveService(testDriveRepo, carRepo)

	bgCtx, stopBackground := context.WithCancel(context.Background())
	defer stopBackground()

	rabbitURL := os.Getenv("RABBITMQ_URL")

	pub, err := broker.NewPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("publisher error: %v", err)
	}
	defer pub.Close()

	go messaging.NewOutboxPoller(outboxRepo, pub).Run(bgCtx)

	resultConsumer, err := messaging.NewOrderResultConsumer(rabbitURL, orderSvc)
	if err != nil {
		log.Fatalf("result consumer error: %v", err)
	}
	defer resultConsumer.Close()
	if err := resultConsumer.Start(bgCtx); err != nil {
		log.Fatalf("result consumer start: %v", err)
	}
	log.Println("outbox poller and result consumer started")

	authHandler := handler.NewAuthHandler(cfg.JWTSecret)
	stockCarHandler := handler.NewStockCarHandler(storageClient)
	carModelHandler := handler.NewCarModelHandler(carModelSvc, detailSvc)
	detailHandler := handler.NewDetailHandler(detailSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	testDriveHandler := handler.NewTestDriveHandler(testDriveSvc)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Group("/auth").Use()
	authHandler.Register(r.Group("/auth"))

	auth := middleware.JWTAuth(cfg.JWTSecret)
	v1 := r.Group("/api/v1")
	v1.Use(auth)

	cars := v1.Group("/cars")
	cars.Use(middleware.RequireRole(middleware.RoleUser, middleware.RoleManager, middleware.RoleAdmin))
	stockCarHandler.Register(cars)

	carModelHandler.Register(v1.Group("/car-models"))
	detailHandler.Register(v1.Group("/details"))
	orderHandler.RegisterInStock(v1.Group("/orders/in-stock"))
	orderHandler.RegisterCustom(v1.Group("/orders/custom"))
	testDriveHandler.Register(v1.Group("/test-drives"))

	srv := server.New(cfg.Port, r)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")
	stopBackground()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
