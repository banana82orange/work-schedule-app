package main

import (
	"fmt"
	"log"
	"net"

	"github.com/portfolio/auth-service/internal/config"
	grpcHandler "github.com/portfolio/auth-service/internal/delivery/grpc"
	"github.com/portfolio/auth-service/internal/infrastructure/repository"
	"github.com/portfolio/auth-service/internal/usecase"
	pb "github.com/portfolio/proto/auth"
	"github.com/portfolio/shared/database"
	"github.com/portfolio/shared/middleware"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()
	fmt.Println(cfg)
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
	userRepo := repository.NewPostgresUserRepository(db)
	roleRepo := repository.NewPostgresRoleRepository(db)
	accessRepo := repository.NewPostgresUserProjectAccessRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, roleRepo, accessRepo, cfg.JWTSecret)
	roleUseCase := usecase.NewRoleUseCase(roleRepo)
	accessUseCase := usecase.NewAccessUseCase(accessRepo)

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(),
			middleware.LoggingInterceptor(),
		),
	)

	// Register auth service
	authServer := grpcHandler.NewAuthServer(authUseCase, roleUseCase, accessUseCase)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	// Start server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Auth service starting on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
