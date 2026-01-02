package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/portfolio/proto/task"
	"github.com/portfolio/shared/database"
	"github.com/portfolio/shared/middleware"
	"github.com/portfolio/task-service/internal/config"
	"github.com/portfolio/task-service/internal/handler"
	"github.com/portfolio/task-service/internal/infrastructure/repository"
	"github.com/portfolio/task-service/internal/usecase"
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
	taskRepo := repository.NewPostgresTaskRepository(db)
	subtaskRepo := repository.NewPostgresSubtaskRepository(db)
	commentRepo := repository.NewPostgresCommentRepository(db)
	attachmentRepo := repository.NewPostgresAttachmentRepository(db)
	tagRepo := repository.NewPostgresTagRepository(db)
	taskTagRepo := repository.NewPostgresTaskTagRepository(db)

	// Initialize use cases
	taskUC := usecase.NewTaskUseCase(taskRepo, subtaskRepo, commentRepo, attachmentRepo, tagRepo, taskTagRepo)
	subtaskUC := usecase.NewSubtaskUseCase(subtaskRepo)
	commentUC := usecase.NewCommentUseCase(commentRepo)
	attachmentUC := usecase.NewAttachmentUseCase(attachmentRepo)
	tagUC := usecase.NewTagUseCase(tagRepo, taskTagRepo)

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(),
			middleware.LoggingInterceptor(),
		),
	)

	// Register task service handler
	taskHandler := handler.NewTaskHandler(taskUC, subtaskUC, commentUC, attachmentUC, tagUC)
	pb.RegisterTaskServiceServer(grpcServer, taskHandler)

	// Start server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Task service starting on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
