package repository

import (
	"context"

	"github.com/portfolio/project-service/internal/domain/entity"
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	Create(ctx context.Context, project *entity.Project) error
	GetByID(ctx context.Context, id int64) (*entity.Project, error)
	Update(ctx context.Context, project *entity.Project) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int, status string) ([]*entity.Project, int, error)
}

// SkillRepository defines the interface for skill data access
type SkillRepository interface {
	Create(ctx context.Context, skill *entity.Skill) error
	GetByID(ctx context.Context, id int64) (*entity.Skill, error)
	GetByName(ctx context.Context, name string) (*entity.Skill, error)
	List(ctx context.Context) ([]*entity.Skill, error)
}

// ProjectSkillRepository defines the interface for project-skill relationship
type ProjectSkillRepository interface {
	Add(ctx context.Context, projectID, skillID int64) error
	Remove(ctx context.Context, projectID, skillID int64) error
	GetByProjectID(ctx context.Context, projectID int64) ([]*entity.Skill, error)
}

// ProjectTechRepository defines the interface for project tech stack
type ProjectTechRepository interface {
	Add(ctx context.Context, projectID int64, techName string) error
	Remove(ctx context.Context, projectID int64, techName string) error
	GetByProjectID(ctx context.Context, projectID int64) ([]string, error)
}

// ProjectImageRepository defines the interface for project images
type ProjectImageRepository interface {
	Add(ctx context.Context, image *entity.ProjectImage) error
	GetByID(ctx context.Context, id int64) (*entity.ProjectImage, error)
	Remove(ctx context.Context, id int64) error
	GetByProjectID(ctx context.Context, projectID int64) ([]*entity.ProjectImage, error)
}

// ProjectLinkRepository defines the interface for project links
type ProjectLinkRepository interface {
	Add(ctx context.Context, link *entity.ProjectLink) error
	GetByID(ctx context.Context, id int64) (*entity.ProjectLink, error)
	Remove(ctx context.Context, id int64) error
	GetByProjectID(ctx context.Context, projectID int64) ([]*entity.ProjectLink, error)
}
