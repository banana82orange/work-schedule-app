package repository

import (
	"context"

	"github.com/portfolio/auth-service/internal/domain/entity"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int) ([]*entity.User, int, error)
}

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id int64) (*entity.Role, error)
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	List(ctx context.Context) ([]*entity.Role, error)
}

// UserProjectAccessRepository defines the interface for user project access data access
type UserProjectAccessRepository interface {
	Set(ctx context.Context, access *entity.UserProjectAccess) error
	Get(ctx context.Context, userID, projectID int64) (*entity.UserProjectAccess, error)
	GetByUserID(ctx context.Context, userID int64) ([]*entity.UserProjectAccess, error)
	GetByProjectID(ctx context.Context, projectID int64) ([]*entity.UserProjectAccess, error)
	Remove(ctx context.Context, userID, projectID int64) error
}
