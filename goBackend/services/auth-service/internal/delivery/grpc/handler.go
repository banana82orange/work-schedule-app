package grpc

import (
	"context"

	"github.com/portfolio/auth-service/internal/domain/entity"
	"github.com/portfolio/auth-service/internal/usecase"
	pb "github.com/portfolio/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthServer implements the AuthService gRPC server
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authUseCase   *usecase.AuthUseCase
	roleUseCase   *usecase.RoleUseCase
	accessUseCase *usecase.AccessUseCase
}

// NewAuthServer creates a new AuthServer
func NewAuthServer(
	authUseCase *usecase.AuthUseCase,
	roleUseCase *usecase.RoleUseCase,
	accessUseCase *usecase.AccessUseCase,
) *AuthServer {
	return &AuthServer{
		authUseCase:   authUseCase,
		roleUseCase:   roleUseCase,
		accessUseCase: accessUseCase,
	}
}

// entityToProto converts entity.User to proto User
func entityToProto(user *entity.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// Register creates a new user
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	role := req.Role
	if role == "" {
		role = "user"
	}

	user, token, err := s.authUseCase.Register(ctx, req.Username, req.Email, req.Password, role)
	if err != nil {
		if err == usecase.ErrUserExists {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		User:  entityToProto(user),
		Token: token,
	}, nil
}

// Login authenticates a user
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, token, err := s.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		User:  entityToProto(user),
		Token: token,
	}, nil
}

// ValidateToken validates a JWT token
func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	user, err := s.authUseCase.ValidateToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid: true,
		User:  entityToProto(user),
	}, nil
}

// GetUser retrieves a user by ID
func (s *AuthServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.authUseCase.GetUser(ctx, req.Id)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{User: entityToProto(user)}, nil
}

// UpdateUser updates a user
func (s *AuthServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, err := s.authUseCase.UpdateUser(ctx, req.Id, req.Username, req.Email, req.Role)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{User: entityToProto(user)}, nil
}

// DeleteUser deletes a user
func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	if err := s.authUseCase.DeleteUser(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

// ListUsers lists users with pagination
func (s *AuthServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, total, err := s.authUseCase.ListUsers(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = entityToProto(user)
	}

	return &pb.ListUsersResponse{
		Users: protoUsers,
		Total: int32(total),
	}, nil
}

// CreateRole creates a new role
func (s *AuthServer) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.RoleResponse, error) {
	role, err := s.roleUseCase.CreateRole(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RoleResponse{
		Role: &pb.Role{
			Id:   role.ID,
			Name: role.Name,
		},
	}, nil
}

// GetRoles lists all roles
func (s *AuthServer) GetRoles(ctx context.Context, req *pb.Empty) (*pb.ListRolesResponse, error) {
	roles, err := s.roleUseCase.ListRoles(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoRoles := make([]*pb.Role, len(roles))
	for i, role := range roles {
		protoRoles[i] = &pb.Role{
			Id:   role.ID,
			Name: role.Name,
		}
	}

	return &pb.ListRolesResponse{Roles: protoRoles}, nil
}

// GetUserProjectAccess gets all project accesses for a user
func (s *AuthServer) GetUserProjectAccess(ctx context.Context, req *pb.GetUserProjectAccessRequest) (*pb.UserProjectAccessResponse, error) {
	accesses, err := s.accessUseCase.GetUserAccess(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoAccesses := make([]*pb.UserProjectAccess, len(accesses))
	for i, access := range accesses {
		protoAccesses[i] = &pb.UserProjectAccess{
			UserId:      access.UserID,
			ProjectId:   access.ProjectID,
			AccessLevel: access.AccessLevel,
		}
	}

	return &pb.UserProjectAccessResponse{Accesses: protoAccesses}, nil
}

// SetUserProjectAccess sets user's access to a project
func (s *AuthServer) SetUserProjectAccess(ctx context.Context, req *pb.SetUserProjectAccessRequest) (*pb.Empty, error) {
	if err := s.accessUseCase.SetAccess(ctx, req.UserId, req.ProjectId, req.AccessLevel); err != nil {
		if err == usecase.ErrInvalidAccessLevel {
			return nil, status.Error(codes.InvalidArgument, "invalid access level")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

// RemoveUserProjectAccess removes user's access to a project
func (s *AuthServer) RemoveUserProjectAccess(ctx context.Context, req *pb.RemoveUserProjectAccessRequest) (*pb.Empty, error) {
	if err := s.accessUseCase.RemoveAccess(ctx, req.UserId, req.ProjectId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}
