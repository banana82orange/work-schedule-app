package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientManager manages gRPC client connections
type ClientManager struct {
	authConn      *grpc.ClientConn
	projectConn   *grpc.ClientConn
	taskConn      *grpc.ClientConn
	analyticsConn *grpc.ClientConn
	mediaConn     *grpc.ClientConn
}

// NewClientManager creates a new ClientManager
func NewClientManager(authURL, projectURL, taskURL, analyticsURL, mediaURL string) (*ClientManager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	// Connect to Auth Service
	authConn, err := grpc.DialContext(ctx, authURL, opts...)
	if err != nil {
		log.Printf("Warning: Failed to connect to Auth service: %v", err)
	}

	// Connect to Project Service
	projectConn, err := grpc.DialContext(ctx, projectURL, opts...)
	if err != nil {
		log.Printf("Warning: Failed to connect to Project service: %v", err)
	}

	// Connect to Task Service
	taskConn, err := grpc.DialContext(ctx, taskURL, opts...)
	if err != nil {
		log.Printf("Warning: Failed to connect to Task service: %v", err)
	}

	// Connect to Analytics Service
	analyticsConn, err := grpc.DialContext(ctx, analyticsURL, opts...)
	if err != nil {
		log.Printf("Warning: Failed to connect to Analytics service: %v", err)
	}

	// Connect to Media Service
	mediaConn, err := grpc.DialContext(ctx, mediaURL, opts...)
	if err != nil {
		log.Printf("Warning: Failed to connect to Media service: %v", err)
	}

	return &ClientManager{
		authConn:      authConn,
		projectConn:   projectConn,
		taskConn:      taskConn,
		analyticsConn: analyticsConn,
		mediaConn:     mediaConn,
	}, nil
}

// GetAuthConn returns the Auth service connection
func (m *ClientManager) GetAuthConn() *grpc.ClientConn {
	return m.authConn
}

// GetProjectConn returns the Project service connection
func (m *ClientManager) GetProjectConn() *grpc.ClientConn {
	return m.projectConn
}

// GetTaskConn returns the Task service connection
func (m *ClientManager) GetTaskConn() *grpc.ClientConn {
	return m.taskConn
}

// GetAnalyticsConn returns the Analytics service connection
func (m *ClientManager) GetAnalyticsConn() *grpc.ClientConn {
	return m.analyticsConn
}

// GetMediaConn returns the Media service connection
func (m *ClientManager) GetMediaConn() *grpc.ClientConn {
	return m.mediaConn
}

// Close closes all connections
func (m *ClientManager) Close() {
	if m.authConn != nil {
		m.authConn.Close()
	}
	if m.projectConn != nil {
		m.projectConn.Close()
	}
	if m.taskConn != nil {
		m.taskConn.Close()
	}
	if m.analyticsConn != nil {
		m.analyticsConn.Close()
	}
	if m.mediaConn != nil {
		m.mediaConn.Close()
	}
}
