package entity

import "time"

// ProjectView represents a project view event
type ProjectView struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	UserID    int64     `json:"user_id"`
	ViewedAt  time.Time `json:"viewed_at"`
}

// NewProjectView creates a new project view
func NewProjectView(projectID, userID int64) *ProjectView {
	return &ProjectView{
		ProjectID: projectID,
		UserID:    userID,
		ViewedAt:  time.Now(),
	}
}

// TaskActivity represents a task activity event
type TaskActivity struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	UserID    int64     `json:"user_id"`
	Action    string    `json:"action"` // created, updated, completed
	CreatedAt time.Time `json:"created_at"`
}

// NewTaskActivity creates a new task activity
func NewTaskActivity(taskID, userID int64, action string) *TaskActivity {
	return &TaskActivity{
		TaskID:    taskID,
		UserID:    userID,
		Action:    action,
		CreatedAt: time.Now(),
	}
}

// Activity action constants
const (
	ActionCreated   = "created"
	ActionUpdated   = "updated"
	ActionCompleted = "completed"
)

// ValidActions returns all valid actions
func ValidActions() []string {
	return []string{ActionCreated, ActionUpdated, ActionCompleted}
}

// ProjectStats represents aggregated project statistics
type ProjectStats struct {
	ProjectID       int64     `json:"project_id"`
	TotalTasks      int     `json:"total_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	ProgressPercent float64   `json:"progress_percent"`
	LastUpdated     time.Time `json:"last_updated"`
}

// NewProjectStats creates a new project stats
func NewProjectStats(projectID int64) *ProjectStats {
	return &ProjectStats{
		ProjectID:       projectID,
		TotalTasks:      0,
		CompletedTasks:  0,
		ProgressPercent: 0,
		LastUpdated:     time.Now(),
	}
}

// UpdateProgress recalculates progress percentage
func (s *ProjectStats) UpdateProgress() {
	if s.TotalTasks > 0 {
		s.ProgressPercent = float64(s.CompletedTasks) / float64(s.TotalTasks) * 100
	} else {
		s.ProgressPercent = 0
	}
	s.LastUpdated = time.Now()
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalProjects  int             `json:"total_projects"`
	ActiveProjects int             `json:"active_projects"`
	TotalTasks     int             `json:"total_tasks"`
	CompletedTasks int             `json:"completed_tasks"`
	PendingTasks   int             `json:"pending_tasks"`
	ProjectStats   []*ProjectStats `json:"project_stats"`
}
