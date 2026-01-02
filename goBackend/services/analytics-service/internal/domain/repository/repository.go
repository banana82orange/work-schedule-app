package repository

import (
	"context"
	"time"

	"github.com/portfolio/analytics-service/internal/domain/entity"
)

// ProjectViewRepository defines the interface for project view data access
type ProjectViewRepository interface {
	Record(ctx context.Context, view *entity.ProjectView) error
	GetByProjectID(ctx context.Context, projectID int64, startDate, endDate *time.Time) ([]*entity.ProjectView, error)
	CountByProjectID(ctx context.Context, projectID int64) (int, error)
}

// TaskActivityRepository defines the interface for task activity data access
type TaskActivityRepository interface {
	Record(ctx context.Context, activity *entity.TaskActivity) error
	GetByTaskID(ctx context.Context, taskID int64) ([]*entity.TaskActivity, error)
	GetByProjectID(ctx context.Context, projectID int64) ([]*entity.TaskActivity, error)
}

// ProjectStatsRepository defines the interface for project stats data access
type ProjectStatsRepository interface {
	Get(ctx context.Context, projectID int64) (*entity.ProjectStats, error)
	Upsert(ctx context.Context, stats *entity.ProjectStats) error
	GetAll(ctx context.Context) ([]*entity.ProjectStats, error)
}
