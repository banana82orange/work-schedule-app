package main

import (
	"fmt"
	"log"
	"net"

	"github.com/portfolio/analytics-service/internal/config"
	grpcHandler "github.com/portfolio/analytics-service/internal/delivery/grpc"
	"github.com/portfolio/analytics-service/internal/infrastructure/repository"
	"github.com/portfolio/analytics-service/internal/usecase"
	"github.com/portfolio/shared/database"
	"github.com/portfolio/shared/middleware"
	"google.golang.org/grpc"
	pb "github.com/portfolio/proto/analytics"
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

	// Initialize repositories
	viewRepo := repository.NewPostgresProjectViewRepository(db)
	actRepo := repository.NewPostgresTaskActivityRepository(db)
	statsRepo := repository.NewPostgresProjectStatsRepository(db)

	// Initialize use cases
	analyticsUseCase := usecase.NewAnalyticsUseCase(viewRepo, actRepo, statsRepo)

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(),
			middleware.LoggingInterceptor(),
		),
	)

	// TODO: Register analytics service handler
	analyticsServer := grpcHandler.NewAnalyticsServer(analyticsUseCase)
	pb.RegisterAnalyticsServiceServer(grpcServer, analyticsServer)

	// Start server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Analytics service starting on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
