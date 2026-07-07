package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/qxblvde/CarService/internal/broker"
	"github.com/qxblvde/CarService/internal/config"
	"github.com/qxblvde/CarService/internal/grpclog"
	"github.com/qxblvde/CarService/internal/middleware"
	"github.com/qxblvde/CarService/internal/server"
	"github.com/qxblvde/CarService/proto/storagepb"
	"github.com/qxblvde/CarService/storage-service/internal/application/service"
	"github.com/qxblvde/CarService/storage-service/internal/infrastructure"
	storagemsg "github.com/qxblvde/CarService/storage-service/internal/infrastructure/messaging"
	"github.com/qxblvde/CarService/storage-service/internal/infrastructure/migrations"
	"github.com/qxblvde/CarService/storage-service/internal/presentation/carserver"
	"github.com/qxblvde/CarService/storage-service/internal/presentation/handler"
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

	rabbitURL := os.Getenv("RABBITMQ_URL")

	pub, err := broker.NewPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("publisher error: %v", err)
	}
	defer pub.Close()

	repo := infrastructure.NewAssemblyOrderRepository(pool)
	svc := service.NewAssemblyOrderService(repo)
	h := handler.NewAssemblyOrderHandler(svc)

	carRepo := infrastructure.NewCarRepository(pool)
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(grpclog.UnaryServer))
	storagepb.RegisterCarStockServer(grpcSrv, carserver.New(carRepo))

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("grpc listen error: %v", err)
	}
	go func() {
		log.Printf("starting grpc server on %s", grpcLis.Addr())
		if err := grpcSrv.Serve(grpcLis); err != nil {
			log.Fatalf("grpc server error: %v", err)
		}
	}()

	consumer, err := storagemsg.NewConsumer(rabbitURL, svc, pub)
	if err != nil {
		log.Fatalf("consumer error: %v", err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("consumer start: %v", err)
	}
	log.Println("consumer started")

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := middleware.JWTAuth(cfg.JWTSecret)
	v1 := r.Group("/api/v1")
	v1.Use(auth)
	ag := v1.Group("/assembly-orders")
	ag.Use(middleware.RequireRole(middleware.RoleWarehouseAdmin))
	h.Register(ag)

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
	cancel()
	grpcSrv.GracefulStop()

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutCancel()
	srv.Shutdown(shutCtx)
}
