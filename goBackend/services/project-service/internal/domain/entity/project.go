package entity

import "time"

// Project represents a project entity
type Project struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	Status      string           `json:"status"`
	Skills      []*Skill         `json:"skills,omitempty"`
	TechStack   []string         `json:"tech_stack,omitempty"`
	Images      []*ProjectImage  `json:"images,omitempty"`
	Links       []*ProjectLink   `json:"links,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// NewProject creates a new project entity
func NewProject(name, description, status string, startDate, endDate *time.Time) *Project {
	now := time.Now()
	if status == "" {
		status = "active"
	}
	return &Project{
		Name:        name,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Skill represents a skill entity
type Skill struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ProjectTech represents project's tech stack
type ProjectTech struct {
	ProjectID int64  `json:"project_id"`
	TechName  string `json:"tech_name"`
}

// ProjectImage represents a project image
type ProjectImage struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

// ProjectLink represents a project link
type ProjectLink struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	LinkURL   string `json:"link_url"`
	LinkType  string `json:"link_type"` // github, live, document
}

// Valid project statuses
const (
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusArchived  = "archived"
	StatusOnHold    = "on_hold"
)

// ValidStatuses returns all valid project statuses
func ValidStatuses() []string {
	return []string{StatusActive, StatusCompleted, StatusArchived, StatusOnHold}
}

// IsValidStatus checks if status is valid
func IsValidStatus(status string) bool {
	for _, s := range ValidStatuses() {
		if s == status {
			return true
		}
	}
	return false
}

// Valid link types
const (
	LinkTypeGitHub   = "github"
	LinkTypeLive     = "live"
	LinkTypeDocument = "document"
)

// ValidLinkTypes returns all valid link types
func ValidLinkTypes() []string {
	return []string{LinkTypeGitHub, LinkTypeLive, LinkTypeDocument}
}
