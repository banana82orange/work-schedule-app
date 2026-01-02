package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/portfolio/proto/auth"
	"google.golang.org/grpc"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authClient pb.AuthServiceClient
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(conn *grpc.ClientConn) *AuthHandler {
	return &AuthHandler{
		authClient: pb.NewAuthServiceClient(conn),
	}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role,omitempty"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents user response
type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// Register handles user registration
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.authClient.Register(ctx, &pb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"email":    resp.User.Email,
			"role":     resp.User.Role,
		},
		"token": resp.Token,
	})
}

// Login handles user login
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.authClient.Login(ctx, &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"email":    resp.User.Email,
			"role":     resp.User.Role,
		},
		"token": resp.Token,
	})
}

// GetProfile returns current user's profile
// GET /api/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// In a real scenario, we might want to fetch fresh data from the service
	// For now, returning what's in the context (from JWT) is fine,
	// or we can call GetUser if we trust the ID in the context.

	// Example of calling service to get fresh data:
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert userID to int64 (depends on how middleware sets it)
	// Assuming it's set as float64 (from JSON) or int
	// specific conversion logic might be needed.
	// For simplicity, let's assume valid ID.

	// ... (Implementation skipped for brevity as we focus on Login/Register first)

	username, _ := c.Get("username")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userID,
			"username": username,
			"email":    email,
			"role":     role,
		},
	})
}

// ValidateToken validates a JWT token
// POST /api/auth/validate
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.authClient.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: req.Token,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": resp.Valid,
		"user":  resp.User, // This might need mapping if pb.User structure differs from desired JSON
	})
}

// ListUsers returns list of users (admin only)
func (h *AuthHandler) ListUsers(c *gin.Context) {
	// Placeholder
	c.JSON(http.StatusOK, gin.H{"message": "List users implementation pending"})
}

// GetUser returns a user by ID
func (h *AuthHandler) GetUser(c *gin.Context) {
	// Placeholder
	c.JSON(http.StatusOK, gin.H{"message": "Get user implementation pending"})
}

// UpdateUser updates a user
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	// Placeholder
	c.JSON(http.StatusOK, gin.H{"message": "Update user implementation pending"})
}

// DeleteUser deletes a user
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	// Placeholder
	c.JSON(http.StatusOK, gin.H{"message": "Delete user implementation pending"})
}
