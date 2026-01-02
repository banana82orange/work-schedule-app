package main

import (
	"fmt"
	"log"
	"net"

	"github.com/portfolio/media-service/internal/config"
	"github.com/portfolio/media-service/internal/infrastructure/repository"
	"github.com/portfolio/media-service/internal/infrastructure/storage"
	"github.com/portfolio/media-service/internal/usecase"
	"github.com/portfolio/shared/database"
	"github.com/portfolio/shared/middleware"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	pool, err := database.NewPool(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	db := pool.GetDB()

	// Initialize storage
	localStorage, err := storage.NewLocalStorage(cfg.StoragePath, cfg.StorageURL)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize repositories
	fileRepo := repository.NewPostgresMediaFileRepository(db)

	// Initialize use cases
	_ = usecase.NewMediaUseCase(fileRepo, localStorage)

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(),
			middleware.LoggingInterceptor(),
		),
	)

	// TODO: Register media service handler

	// Start server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Media service starting on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
