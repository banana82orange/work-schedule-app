package main

import (
	"fmt"
	"log"
	"net"

	"github.com/portfolio/project-service/internal/config"
	"github.com/portfolio/project-service/internal/handler"
	"github.com/portfolio/project-service/internal/infrastructure/repository"
	"github.com/portfolio/project-service/internal/usecase"
	pb "github.com/portfolio/proto/project"
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

	// Initialize repositories
	projectRepo := repository.NewPostgresProjectRepository(db)
	skillRepo := repository.NewPostgresSkillRepository(db)
	projectSkillRepo := repository.NewPostgresProjectSkillRepository(db)
	techRepo := repository.NewPostgresProjectTechRepository(db)
	imageRepo := repository.NewPostgresProjectImageRepository(db)
	linkRepo := repository.NewPostgresProjectLinkRepository(db)

	// Initialize use cases
	projectUC := usecase.NewProjectUseCase(projectRepo, skillRepo, projectSkillRepo, techRepo, imageRepo, linkRepo)
	skillUC := usecase.NewSkillUseCase(skillRepo)
	projectSkillUC := usecase.NewProjectSkillUseCase(projectSkillRepo)
	techUC := usecase.NewTechUseCase(techRepo)
	imageUC := usecase.NewImageUseCase(imageRepo)
	linkUC := usecase.NewLinkUseCase(linkRepo)

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(),
			middleware.LoggingInterceptor(),
		),
	)

	// Register project service handler
	projectHandler := handler.NewProjectHandler(projectUC, skillUC, projectSkillUC, techUC, imageUC, linkUC)
	pb.RegisterProjectServiceServer(grpcServer, projectHandler)

	// Start server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Project service starting on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
