package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/portfolio/project-service/internal/domain/entity"
	"github.com/portfolio/project-service/internal/domain/repository"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrSkillNotFound   = errors.New("skill not found")
	ErrImageNotFound   = errors.New("image not found")
	ErrLinkNotFound    = errors.New("link not found")
)

// ProjectUseCase handles project business logic
type ProjectUseCase struct {
	projectRepo      repository.ProjectRepository
	skillRepo        repository.SkillRepository
	projectSkillRepo repository.ProjectSkillRepository
	techRepo         repository.ProjectTechRepository
	imageRepo        repository.ProjectImageRepository
	linkRepo         repository.ProjectLinkRepository
}

// NewProjectUseCase creates a new ProjectUseCase
func NewProjectUseCase(
	projectRepo repository.ProjectRepository,
	skillRepo repository.SkillRepository,
	projectSkillRepo repository.ProjectSkillRepository,
	techRepo repository.ProjectTechRepository,
	imageRepo repository.ProjectImageRepository,
	linkRepo repository.ProjectLinkRepository,
) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo:      projectRepo,
		skillRepo:        skillRepo,
		projectSkillRepo: projectSkillRepo,
		techRepo:         techRepo,
		imageRepo:        imageRepo,
		linkRepo:         linkRepo,
	}
}

// CreateProject creates a new project
func (uc *ProjectUseCase) CreateProject(ctx context.Context, name, description, status string, startDate, endDate *time.Time) (*entity.Project, error) {
	project := entity.NewProject(name, description, status, startDate, endDate)
	if err := uc.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

// GetProject retrieves a project by ID with all related data
func (uc *ProjectUseCase) GetProject(ctx context.Context, id int64) (*entity.Project, error) {
	project, err := uc.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProjectNotFound
	}

	// Load related data
	skills, _ := uc.projectSkillRepo.GetByProjectID(ctx, id)
	project.Skills = skills

	techStack, _ := uc.techRepo.GetByProjectID(ctx, id)
	project.TechStack = techStack

	images, _ := uc.imageRepo.GetByProjectID(ctx, id)
	project.Images = images

	links, _ := uc.linkRepo.GetByProjectID(ctx, id)
	project.Links = links

	return project, nil
}

// UpdateProject updates a project
func (uc *ProjectUseCase) UpdateProject(ctx context.Context, id int64, name, description, status string, startDate, endDate *time.Time) (*entity.Project, error) {
	project, err := uc.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProjectNotFound
	}

	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}
	if status != "" {
		project.Status = status
	}
	if startDate != nil {
		project.StartDate = startDate
	}
	if endDate != nil {
		project.EndDate = endDate
	}
	project.UpdatedAt = time.Now()

	if err := uc.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return uc.GetProject(ctx, id)
}

// DeleteProject deletes a project
func (uc *ProjectUseCase) DeleteProject(ctx context.Context, id int64) error {
	return uc.projectRepo.Delete(ctx, id)
}

// ListProjects lists projects with pagination
func (uc *ProjectUseCase) ListProjects(ctx context.Context, page, limit int, status string) ([]*entity.Project, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.projectRepo.List(ctx, page, limit, status)
}

// SkillUseCase handles skill business logic
type SkillUseCase struct {
	skillRepo repository.SkillRepository
}

// NewSkillUseCase creates a new SkillUseCase
func NewSkillUseCase(skillRepo repository.SkillRepository) *SkillUseCase {
	return &SkillUseCase{skillRepo: skillRepo}
}

// CreateSkill creates a new skill
func (uc *SkillUseCase) CreateSkill(ctx context.Context, name string) (*entity.Skill, error) {
	skill := &entity.Skill{Name: name}
	if err := uc.skillRepo.Create(ctx, skill); err != nil {
		return nil, err
	}
	return skill, nil
}

// ListSkills lists all skills
func (uc *SkillUseCase) ListSkills(ctx context.Context) ([]*entity.Skill, error) {
	return uc.skillRepo.List(ctx)
}

// ProjectSkillUseCase handles project-skill relationships
type ProjectSkillUseCase struct {
	projectSkillRepo repository.ProjectSkillRepository
}

// NewProjectSkillUseCase creates a new ProjectSkillUseCase
func NewProjectSkillUseCase(projectSkillRepo repository.ProjectSkillRepository) *ProjectSkillUseCase {
	return &ProjectSkillUseCase{projectSkillRepo: projectSkillRepo}
}

// AddSkill adds a skill to a project
func (uc *ProjectSkillUseCase) AddSkill(ctx context.Context, projectID, skillID int64) error {
	return uc.projectSkillRepo.Add(ctx, projectID, skillID)
}

// RemoveSkill removes a skill from a project
func (uc *ProjectSkillUseCase) RemoveSkill(ctx context.Context, projectID, skillID int64) error {
	return uc.projectSkillRepo.Remove(ctx, projectID, skillID)
}

// TechUseCase handles project tech stack
type TechUseCase struct {
	techRepo repository.ProjectTechRepository
}

// NewTechUseCase creates a new TechUseCase
func NewTechUseCase(techRepo repository.ProjectTechRepository) *TechUseCase {
	return &TechUseCase{techRepo: techRepo}
}

// AddTech adds a technology to a project
func (uc *TechUseCase) AddTech(ctx context.Context, projectID int64, techName string) error {
	return uc.techRepo.Add(ctx, projectID, techName)
}

// RemoveTech removes a technology from a project
func (uc *TechUseCase) RemoveTech(ctx context.Context, projectID int64, techName string) error {
	return uc.techRepo.Remove(ctx, projectID, techName)
}

// ImageUseCase handles project images
type ImageUseCase struct {
	imageRepo repository.ProjectImageRepository
}

// NewImageUseCase creates a new ImageUseCase
func NewImageUseCase(imageRepo repository.ProjectImageRepository) *ImageUseCase {
	return &ImageUseCase{imageRepo: imageRepo}
}

// AddImage adds an image to a project
func (uc *ImageUseCase) AddImage(ctx context.Context, projectID int64, imageURL, description string) (*entity.ProjectImage, error) {
	image := &entity.ProjectImage{
		ProjectID:   projectID,
		ImageURL:    imageURL,
		Description: description,
		UploadedAt:  time.Now(),
	}
	if err := uc.imageRepo.Add(ctx, image); err != nil {
		return nil, err
	}
	return image, nil
}

// RemoveImage removes an image
func (uc *ImageUseCase) RemoveImage(ctx context.Context, id int64) error {
	return uc.imageRepo.Remove(ctx, id)
}

// GetImages gets all images for a project
func (uc *ImageUseCase) GetImages(ctx context.Context, projectID int64) ([]*entity.ProjectImage, error) {
	return uc.imageRepo.GetByProjectID(ctx, projectID)
}

// LinkUseCase handles project links
type LinkUseCase struct {
	linkRepo repository.ProjectLinkRepository
}

// NewLinkUseCase creates a new LinkUseCase
func NewLinkUseCase(linkRepo repository.ProjectLinkRepository) *LinkUseCase {
	return &LinkUseCase{linkRepo: linkRepo}
}

// AddLink adds a link to a project
func (uc *LinkUseCase) AddLink(ctx context.Context, projectID int64, linkURL, linkType string) (*entity.ProjectLink, error) {
	link := &entity.ProjectLink{
		ProjectID: projectID,
		LinkURL:   linkURL,
		LinkType:  linkType,
	}
	if err := uc.linkRepo.Add(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

// RemoveLink removes a link
func (uc *LinkUseCase) RemoveLink(ctx context.Context, id int64) error {
	return uc.linkRepo.Remove(ctx, id)
}

// GetLinks gets all links for a project
func (uc *LinkUseCase) GetLinks(ctx context.Context, projectID int64) ([]*entity.ProjectLink, error) {
	return uc.linkRepo.GetByProjectID(ctx, projectID)
}
