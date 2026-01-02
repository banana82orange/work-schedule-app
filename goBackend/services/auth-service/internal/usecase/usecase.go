package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/portfolio/auth-service/internal/domain/entity"
	"github.com/portfolio/auth-service/internal/domain/repository"
	"github.com/portfolio/shared/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidAccessLevel = errors.New("invalid access level")
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo    repository.UserRepository
	roleRepo    repository.RoleRepository
	accessRepo  repository.UserProjectAccessRepository
	tokenSvc    *jwt.TokenService
}

// NewAuthUseCase creates a new AuthUseCase
func NewAuthUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	accessRepo repository.UserProjectAccessRepository,
	jwtSecret string,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		accessRepo: accessRepo,
		tokenSvc:   jwt.NewTokenService(jwtSecret, 24*time.Hour),
	}
}

// Register creates a new user
func (uc *AuthUseCase) Register(ctx context.Context, username, email, password, role string) (*entity.User, string, error) {
	// Check if user exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, "", ErrUserExists
	}

	existingUser, _ = uc.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, "", ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := entity.NewUser(username, email, string(hashedPassword), role)
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := uc.tokenSvc.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login authenticates a user
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*entity.User, string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := uc.tokenSvc.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ValidateToken validates a JWT token
func (uc *AuthUseCase) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	claims, err := uc.tokenSvc.ValidateToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (uc *AuthUseCase) GetUser(ctx context.Context, id int64) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser updates a user
func (uc *AuthUseCase) UpdateUser(ctx context.Context, id int64, username, email, role string) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if role != "" {
		user.Role = role
	}
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (uc *AuthUseCase) DeleteUser(ctx context.Context, id int64) error {
	return uc.userRepo.Delete(ctx, id)
}

// ListUsers lists users with pagination
func (uc *AuthUseCase) ListUsers(ctx context.Context, page, limit int) ([]*entity.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.userRepo.List(ctx, page, limit)
}

// RoleUseCase handles role business logic
type RoleUseCase struct {
	roleRepo repository.RoleRepository
}

// NewRoleUseCase creates a new RoleUseCase
func NewRoleUseCase(roleRepo repository.RoleRepository) *RoleUseCase {
	return &RoleUseCase{roleRepo: roleRepo}
}

// CreateRole creates a new role
func (uc *RoleUseCase) CreateRole(ctx context.Context, name string) (*entity.Role, error) {
	role := &entity.Role{Name: name}
	if err := uc.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

// ListRoles lists all roles
func (uc *RoleUseCase) ListRoles(ctx context.Context) ([]*entity.Role, error) {
	return uc.roleRepo.List(ctx)
}

// AccessUseCase handles project access business logic
type AccessUseCase struct {
	accessRepo repository.UserProjectAccessRepository
}

// NewAccessUseCase creates a new AccessUseCase
func NewAccessUseCase(accessRepo repository.UserProjectAccessRepository) *AccessUseCase {
	return &AccessUseCase{accessRepo: accessRepo}
}

// SetAccess sets user's access to a project
func (uc *AccessUseCase) SetAccess(ctx context.Context, userID, projectID int64, accessLevel string) error {
	if !entity.IsValidAccessLevel(accessLevel) {
		return ErrInvalidAccessLevel
	}

	access := &entity.UserProjectAccess{
		UserID:      userID,
		ProjectID:   projectID,
		AccessLevel: accessLevel,
	}
	return uc.accessRepo.Set(ctx, access)
}

// GetUserAccess gets all project accesses for a user
func (uc *AccessUseCase) GetUserAccess(ctx context.Context, userID int64) ([]*entity.UserProjectAccess, error) {
	return uc.accessRepo.GetByUserID(ctx, userID)
}

// RemoveAccess removes user's access to a project
func (uc *AccessUseCase) RemoveAccess(ctx context.Context, userID, projectID int64) error {
	return uc.accessRepo.Remove(ctx, userID, projectID)
}
