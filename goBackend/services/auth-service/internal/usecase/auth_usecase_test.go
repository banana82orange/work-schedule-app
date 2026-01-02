package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/portfolio/auth-service/internal/domain/entity"
)

// MockUserRepository is a manual mock
type MockUserRepository struct {
	users map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	// Simulate DB ID generation
	user.ID = int64(len(m.users) + 1)
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// Implement other methods as no-ops or panics if not used in tested paths
func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) { return nil, nil }
func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error { return nil }
func (m *MockUserRepository) Delete(ctx context.Context, id int64) error { return nil }
func (m *MockUserRepository) List(ctx context.Context, page, limit int) ([]*entity.User, int, error) { return nil, 0, nil }


func TestAuthUseCase_Register(t *testing.T) {
	mockRepo := NewMockUserRepository()
	// Mock other repos as nil since they aren't used in Register (except maybe for checking, but logic shows only userRepo is critical for this path if we ignore role/access for now or mock them if strictly needed.
	// actually Register uses: userRepo.GetByEmail, userRepo.GetByUsername, userRepo.Create.
	// It relies on tokenSvc internally.

	uc := NewAuthUseCase(mockRepo, nil, nil, "secret")

	tests := []struct {
		name    string
		username string
		email   string
		password string
		role    string
		wantErr bool
	}{
		{
			name:    "Success",
			username: "testuser",
			email:   "test@example.com",
			password: "password123",
			role:    "user",
			wantErr: false,
		},
		{
			name:    "Duplicate Email",
			username: "otheruser",
			email:   "test@example.com", // Same as above
			password: "password123",
			role:    "user",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, token, err := uc.Register(context.Background(), tt.username, tt.email, tt.password, tt.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUseCase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if user == nil {
					t.Error("AuthUseCase.Register() user should not be nil")
				}
				if token == "" {
					t.Error("AuthUseCase.Register() token should not be empty")
				}
				if user.Email != tt.email {
					t.Errorf("AuthUseCase.Register() user email = %v, want %v", user.Email, tt.email)
				}
			}
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := NewAuthUseCase(mockRepo, nil, nil, "secret")

	// Pre-seed a user
	uc.Register(context.Background(), "loginuser", "login@example.com", "password123", "user")

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Success",
			email:    "login@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Invalid Password",
			email:    "login@example.com",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "User Not Found",
			email:    "notfound@example.com",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, token, err := uc.Login(context.Background(), tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUseCase.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if user == nil {
					t.Error("AuthUseCase.Login() user should not be nil")
				}
				if token == "" {
					t.Error("AuthUseCase.Login() token should not be empty")
				}
			}
		})
	}
}
