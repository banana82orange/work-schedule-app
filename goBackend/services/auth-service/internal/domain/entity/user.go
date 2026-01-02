package entity

import "time"

// User represents a user entity
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser creates a new user entity
func NewUser(username, email, passwordHash, role string) *User {
	now := time.Now()
	if role == "" {
		role = "user"
	}
	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Role represents a role entity
type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UserProjectAccess represents user's access to a project
type UserProjectAccess struct {
	UserID      int64  `json:"user_id"`
	ProjectID   int64  `json:"project_id"`
	AccessLevel string `json:"access_level"` // read, write, admin
}

// AccessLevel constants
const (
	AccessLevelRead  = "read"
	AccessLevelWrite = "write"
	AccessLevelAdmin = "admin"
)

// ValidAccessLevels returns all valid access levels
func ValidAccessLevels() []string {
	return []string{AccessLevelRead, AccessLevelWrite, AccessLevelAdmin}
}

// IsValidAccessLevel checks if access level is valid
func IsValidAccessLevel(level string) bool {
	for _, valid := range ValidAccessLevels() {
		if valid == level {
			return true
		}
	}
	return false
}

// HasWriteAccess checks if user has write access
func (a *UserProjectAccess) HasWriteAccess() bool {
	return a.AccessLevel == AccessLevelWrite || a.AccessLevel == AccessLevelAdmin
}

// HasAdminAccess checks if user has admin access
func (a *UserProjectAccess) HasAdminAccess() bool {
	return a.AccessLevel == AccessLevelAdmin
}
