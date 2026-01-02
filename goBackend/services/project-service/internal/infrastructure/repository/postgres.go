package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/portfolio/project-service/internal/domain/entity"
)

// PostgresProjectRepository implements ProjectRepository
type PostgresProjectRepository struct {
	db *sql.DB
}

// NewPostgresProjectRepository creates a new PostgresProjectRepository
func NewPostgresProjectRepository(db *sql.DB) *PostgresProjectRepository {
	return &PostgresProjectRepository{db: db}
}

// Create creates a new project
func (r *PostgresProjectRepository) Create(ctx context.Context, project *entity.Project) error {
	query := `
		INSERT INTO projects (name, description, start_date, end_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	return r.db.QueryRowContext(
		ctx, query,
		project.Name, project.Description, project.StartDate, project.EndDate,
		project.Status, project.CreatedAt, project.UpdatedAt,
	).Scan(&project.ID)
}

// GetByID gets a project by ID
func (r *PostgresProjectRepository) GetByID(ctx context.Context, id int64) (*entity.Project, error) {
	query := `
		SELECT id, name, description, start_date, end_date, status, created_at, updated_at
		FROM projects WHERE id = $1
	`
	project := &entity.Project{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&project.ID, &project.Name, &project.Description,
		&project.StartDate, &project.EndDate, &project.Status,
		&project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// Update updates a project
func (r *PostgresProjectRepository) Update(ctx context.Context, project *entity.Project) error {
	query := `
		UPDATE projects SET name = $1, description = $2, start_date = $3,
		end_date = $4, status = $5, updated_at = $6 WHERE id = $7
	`
	project.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		project.Name, project.Description, project.StartDate,
		project.EndDate, project.Status, project.UpdatedAt, project.ID,
	)
	return err
}

// Delete deletes a project
func (r *PostgresProjectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List lists projects with pagination
func (r *PostgresProjectRepository) List(ctx context.Context, page, limit int, status string) ([]*entity.Project, int, error) {
	offset := (page - 1) * limit

	// Build query based on status filter
	var countQuery, query string
	var args []interface{}

	if status != "" {
		countQuery = `SELECT COUNT(*) FROM projects WHERE status = $1`
		query = `
			SELECT id, name, description, start_date, end_date, status, created_at, updated_at
			FROM projects WHERE status = $1 ORDER BY id LIMIT $2 OFFSET $3
		`
		args = []interface{}{status, limit, offset}
	} else {
		countQuery = `SELECT COUNT(*) FROM projects`
		query = `
			SELECT id, name, description, start_date, end_date, status, created_at, updated_at
			FROM projects ORDER BY id LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	// Get total count
	var total int
	if status != "" {
		if err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	// Get projects
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []*entity.Project
	for rows.Next() {
		project := &entity.Project{}
		if err := rows.Scan(
			&project.ID, &project.Name, &project.Description,
			&project.StartDate, &project.EndDate, &project.Status,
			&project.CreatedAt, &project.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		projects = append(projects, project)
	}

	return projects, total, nil
}

// PostgresSkillRepository implements SkillRepository
type PostgresSkillRepository struct {
	db *sql.DB
}

// NewPostgresSkillRepository creates a new PostgresSkillRepository
func NewPostgresSkillRepository(db *sql.DB) *PostgresSkillRepository {
	return &PostgresSkillRepository{db: db}
}

// Create creates a new skill
func (r *PostgresSkillRepository) Create(ctx context.Context, skill *entity.Skill) error {
	query := `INSERT INTO skills (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, skill.Name).Scan(&skill.ID)
}

// GetByID gets a skill by ID
func (r *PostgresSkillRepository) GetByID(ctx context.Context, id int64) (*entity.Skill, error) {
	query := `SELECT id, name FROM skills WHERE id = $1`
	skill := &entity.Skill{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&skill.ID, &skill.Name)
	if err != nil {
		return nil, err
	}
	return skill, nil
}

// GetByName gets a skill by name
func (r *PostgresSkillRepository) GetByName(ctx context.Context, name string) (*entity.Skill, error) {
	query := `SELECT id, name FROM skills WHERE name = $1`
	skill := &entity.Skill{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&skill.ID, &skill.Name)
	if err != nil {
		return nil, err
	}
	return skill, nil
}

// List lists all skills
func (r *PostgresSkillRepository) List(ctx context.Context) ([]*entity.Skill, error) {
	query := `SELECT id, name FROM skills ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*entity.Skill
	for rows.Next() {
		skill := &entity.Skill{}
		if err := rows.Scan(&skill.ID, &skill.Name); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	return skills, nil
}

// PostgresProjectSkillRepository implements ProjectSkillRepository
type PostgresProjectSkillRepository struct {
	db *sql.DB
}

// NewPostgresProjectSkillRepository creates a new repository
func NewPostgresProjectSkillRepository(db *sql.DB) *PostgresProjectSkillRepository {
	return &PostgresProjectSkillRepository{db: db}
}

// Add adds a skill to a project
func (r *PostgresProjectSkillRepository) Add(ctx context.Context, projectID, skillID int64) error {
	query := `INSERT INTO project_skills (project_id, skill_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, projectID, skillID)
	return err
}

// Remove removes a skill from a project
func (r *PostgresProjectSkillRepository) Remove(ctx context.Context, projectID, skillID int64) error {
	query := `DELETE FROM project_skills WHERE project_id = $1 AND skill_id = $2`
	_, err := r.db.ExecContext(ctx, query, projectID, skillID)
	return err
}

// GetByProjectID gets all skills for a project
func (r *PostgresProjectSkillRepository) GetByProjectID(ctx context.Context, projectID int64) ([]*entity.Skill, error) {
	query := `
		SELECT s.id, s.name FROM skills s
		INNER JOIN project_skills ps ON s.id = ps.skill_id
		WHERE ps.project_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*entity.Skill
	for rows.Next() {
		skill := &entity.Skill{}
		if err := rows.Scan(&skill.ID, &skill.Name); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	return skills, nil
}

// PostgresProjectTechRepository implements ProjectTechRepository
type PostgresProjectTechRepository struct {
	db *sql.DB
}

// NewPostgresProjectTechRepository creates a new repository
func NewPostgresProjectTechRepository(db *sql.DB) *PostgresProjectTechRepository {
	return &PostgresProjectTechRepository{db: db}
}

// Add adds a technology to a project
func (r *PostgresProjectTechRepository) Add(ctx context.Context, projectID int64, techName string) error {
	query := `INSERT INTO project_tech (project_id, tech_name) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, projectID, techName)
	return err
}

// Remove removes a technology from a project
func (r *PostgresProjectTechRepository) Remove(ctx context.Context, projectID int64, techName string) error {
	query := `DELETE FROM project_tech WHERE project_id = $1 AND tech_name = $2`
	_, err := r.db.ExecContext(ctx, query, projectID, techName)
	return err
}

// GetByProjectID gets all technologies for a project
func (r *PostgresProjectTechRepository) GetByProjectID(ctx context.Context, projectID int64) ([]string, error) {
	query := `SELECT tech_name FROM project_tech WHERE project_id = $1`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var techs []string
	for rows.Next() {
		var tech string
		if err := rows.Scan(&tech); err != nil {
			return nil, err
		}
		techs = append(techs, tech)
	}
	return techs, nil
}

// PostgresProjectImageRepository implements ProjectImageRepository
type PostgresProjectImageRepository struct {
	db *sql.DB
}

// NewPostgresProjectImageRepository creates a new repository
func NewPostgresProjectImageRepository(db *sql.DB) *PostgresProjectImageRepository {
	return &PostgresProjectImageRepository{db: db}
}

// Add adds an image to a project
func (r *PostgresProjectImageRepository) Add(ctx context.Context, image *entity.ProjectImage) error {
	query := `
		INSERT INTO project_images (project_id, image_url, description, uploaded_at)
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		image.ProjectID, image.ImageURL, image.Description, image.UploadedAt,
	).Scan(&image.ID)
}

// GetByID gets an image by ID
func (r *PostgresProjectImageRepository) GetByID(ctx context.Context, id int64) (*entity.ProjectImage, error) {
	query := `SELECT id, project_id, image_url, description, uploaded_at FROM project_images WHERE id = $1`
	image := &entity.ProjectImage{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&image.ID, &image.ProjectID, &image.ImageURL, &image.Description, &image.UploadedAt,
	)
	if err != nil {
		return nil, err
	}
	return image, nil
}

// Remove removes an image
func (r *PostgresProjectImageRepository) Remove(ctx context.Context, id int64) error {
	query := `DELETE FROM project_images WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByProjectID gets all images for a project
func (r *PostgresProjectImageRepository) GetByProjectID(ctx context.Context, projectID int64) ([]*entity.ProjectImage, error) {
	query := `SELECT id, project_id, image_url, description, uploaded_at FROM project_images WHERE project_id = $1`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*entity.ProjectImage
	for rows.Next() {
		image := &entity.ProjectImage{}
		if err := rows.Scan(&image.ID, &image.ProjectID, &image.ImageURL, &image.Description, &image.UploadedAt); err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

// PostgresProjectLinkRepository implements ProjectLinkRepository
type PostgresProjectLinkRepository struct {
	db *sql.DB
}

// NewPostgresProjectLinkRepository creates a new repository
func NewPostgresProjectLinkRepository(db *sql.DB) *PostgresProjectLinkRepository {
	return &PostgresProjectLinkRepository{db: db}
}

// Add adds a link to a project
func (r *PostgresProjectLinkRepository) Add(ctx context.Context, link *entity.ProjectLink) error {
	query := `
		INSERT INTO project_links (project_id, link_url, link_type)
		VALUES ($1, $2, $3) RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, link.ProjectID, link.LinkURL, link.LinkType).Scan(&link.ID)
}

// GetByID gets a link by ID
func (r *PostgresProjectLinkRepository) GetByID(ctx context.Context, id int64) (*entity.ProjectLink, error) {
	query := `SELECT id, project_id, link_url, link_type FROM project_links WHERE id = $1`
	link := &entity.ProjectLink{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&link.ID, &link.ProjectID, &link.LinkURL, &link.LinkType)
	if err != nil {
		return nil, err
	}
	return link, nil
}

// Remove removes a link
func (r *PostgresProjectLinkRepository) Remove(ctx context.Context, id int64) error {
	query := `DELETE FROM project_links WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByProjectID gets all links for a project
func (r *PostgresProjectLinkRepository) GetByProjectID(ctx context.Context, projectID int64) ([]*entity.ProjectLink, error) {
	query := `SELECT id, project_id, link_url, link_type FROM project_links WHERE project_id = $1`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*entity.ProjectLink
	for rows.Next() {
		link := &entity.ProjectLink{}
		if err := rows.Scan(&link.ID, &link.ProjectID, &link.LinkURL, &link.LinkType); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
