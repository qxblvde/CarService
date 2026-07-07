package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port      int
	DBURL     string
	JWTSecret string

	GRPCPort        int
	StorageGRPCAddr string
}

func Load() (Config, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return Config{}, fmt.Errorf("DB_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	port := 8080
	if raw := os.Getenv("PORT"); raw != "" {
		p, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("PORT must be a number: %w", err)
		}
		port = p
	}

	grpcPort := 9091
	if raw := os.Getenv("GRPC_PORT"); raw != "" {
		p, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("GRPC_PORT must be a number: %w", err)
		}
		grpcPort = p
	}

	storageGRPCAddr := os.Getenv("STORAGE_GRPC_ADDR")
	if storageGRPCAddr == "" {
		storageGRPCAddr = "localhost:9091"
	}

	return Config{
		Port:            port,
		DBURL:           dbURL,
		JWTSecret:       jwtSecret,
		GRPCPort:        grpcPort,
		StorageGRPCAddr: storageGRPCAddr,
	}, nil
}
