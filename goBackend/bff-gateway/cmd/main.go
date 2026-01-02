package main

import (
	"fmt"
	"log"

	"github.com/portfolio/bff-gateway/internal/config"
	"github.com/portfolio/bff-gateway/internal/grpc"
	"github.com/portfolio/bff-gateway/internal/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize gRPC clients
	clientManager, err := grpc.NewClientManager(
		cfg.AuthServiceURL,
		cfg.ProjectServiceURL,
		cfg.TaskServiceURL,
		cfg.AnalyticsServiceURL,
		cfg.MediaServiceURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer clientManager.Close()

	// Setup router
	r := router.SetupRouter(cfg.JWTSecret, clientManager)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	log.Printf("BFF Gateway starting on %s", addr)
	log.Printf("Service URLs:")
	log.Printf("  Auth:      %s", cfg.AuthServiceURL)
	log.Printf("  Project:   %s", cfg.ProjectServiceURL)
	log.Printf("  Task:      %s", cfg.TaskServiceURL)
	log.Printf("  Analytics: %s", cfg.AnalyticsServiceURL)
	log.Printf("  Media:     %s", cfg.MediaServiceURL)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
