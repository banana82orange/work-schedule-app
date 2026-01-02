package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/portfolio/auth-service/internal/domain/entity"
)

// PostgresUserRepository implements UserRepository
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRowContext(
		ctx, query,
		user.Username, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
}

// GetByID gets a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByEmail gets a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByUsername gets a user by username
func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users WHERE username = $1
	`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates a user
func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users SET username = $1, email = $2, role = $3, updated_at = $4
		WHERE id = $5
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Role, user.UpdatedAt, user.ID)
	return err
}

// Delete deletes a user
func (r *PostgresUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List lists users with pagination
func (r *PostgresUserRepository) List(ctx context.Context, page, limit int) ([]*entity.User, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get users
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users ORDER BY id LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Role, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// PostgresRoleRepository implements RoleRepository
type PostgresRoleRepository struct {
	db *sql.DB
}

// NewPostgresRoleRepository creates a new PostgresRoleRepository
func NewPostgresRoleRepository(db *sql.DB) *PostgresRoleRepository {
	return &PostgresRoleRepository{db: db}
}

// Create creates a new role
func (r *PostgresRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	query := `INSERT INTO roles (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, role.Name).Scan(&role.ID)
}

// GetByID gets a role by ID
func (r *PostgresRoleRepository) GetByID(ctx context.Context, id int64) (*entity.Role, error) {
	query := `SELECT id, name FROM roles WHERE id = $1`
	role := &entity.Role{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// GetByName gets a role by name
func (r *PostgresRoleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	query := `SELECT id, name FROM roles WHERE name = $1`
	role := &entity.Role{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// List lists all roles
func (r *PostgresRoleRepository) List(ctx context.Context) ([]*entity.Role, error) {
	query := `SELECT id, name FROM roles ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		role := &entity.Role{}
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// PostgresUserProjectAccessRepository implements UserProjectAccessRepository
type PostgresUserProjectAccessRepository struct {
	db *sql.DB
}

// NewPostgresUserProjectAccessRepository creates a new PostgresUserProjectAccessRepository
func NewPostgresUserProjectAccessRepository(db *sql.DB) *PostgresUserProjectAccessRepository {
	return &PostgresUserProjectAccessRepository{db: db}
}

// Set sets user's access to a project (upsert)
func (r *PostgresUserProjectAccessRepository) Set(ctx context.Context, access *entity.UserProjectAccess) error {
	query := `
		INSERT INTO user_project_access (user_id, project_id, access_level)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, project_id) DO UPDATE SET access_level = $3
	`
	_, err := r.db.ExecContext(ctx, query, access.UserID, access.ProjectID, access.AccessLevel)
	return err
}

// Get gets user's access to a specific project
func (r *PostgresUserProjectAccessRepository) Get(ctx context.Context, userID, projectID int64) (*entity.UserProjectAccess, error) {
	query := `
		SELECT user_id, project_id, access_level
		FROM user_project_access WHERE user_id = $1 AND project_id = $2
	`
	access := &entity.UserProjectAccess{}
	err := r.db.QueryRowContext(ctx, query, userID, projectID).Scan(
		&access.UserID, &access.ProjectID, &access.AccessLevel,
	)
	if err != nil {
		return nil, err
	}
	return access, nil
}

// GetByUserID gets all project accesses for a user
func (r *PostgresUserProjectAccessRepository) GetByUserID(ctx context.Context, userID int64) ([]*entity.UserProjectAccess, error) {
	query := `
		SELECT user_id, project_id, access_level
		FROM user_project_access WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accesses []*entity.UserProjectAccess
	for rows.Next() {
		access := &entity.UserProjectAccess{}
		if err := rows.Scan(&access.UserID, &access.ProjectID, &access.AccessLevel); err != nil {
			return nil, err
		}
		accesses = append(accesses, access)
	}
	return accesses, nil
}

// GetByProjectID gets all user accesses for a project
func (r *PostgresUserProjectAccessRepository) GetByProjectID(ctx context.Context, projectID int64) ([]*entity.UserProjectAccess, error) {
	query := `
		SELECT user_id, project_id, access_level
		FROM user_project_access WHERE project_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accesses []*entity.UserProjectAccess
	for rows.Next() {
		access := &entity.UserProjectAccess{}
		if err := rows.Scan(&access.UserID, &access.ProjectID, &access.AccessLevel); err != nil {
			return nil, err
		}
		accesses = append(accesses, access)
	}
	return accesses, nil
}

// Remove removes user's access to a project
func (r *PostgresUserProjectAccessRepository) Remove(ctx context.Context, userID, projectID int64) error {
	query := `DELETE FROM user_project_access WHERE user_id = $1 AND project_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, projectID)
	return err
}
