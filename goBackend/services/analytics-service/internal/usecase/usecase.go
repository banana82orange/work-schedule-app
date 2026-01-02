package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/portfolio/analytics-service/internal/domain/entity"
	"github.com/portfolio/analytics-service/internal/domain/repository"
)

var (
	ErrProjectStatsNotFound = errors.New("project stats not found")
)

// AnalyticsUseCase handles analytics business logic
type AnalyticsUseCase struct {
	viewRepo  repository.ProjectViewRepository
	actRepo   repository.TaskActivityRepository
	statsRepo repository.ProjectStatsRepository
}

// NewAnalyticsUseCase creates a new AnalyticsUseCase
func NewAnalyticsUseCase(
	viewRepo repository.ProjectViewRepository,
	actRepo repository.TaskActivityRepository,
	statsRepo repository.ProjectStatsRepository,
) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		viewRepo:  viewRepo,
		actRepo:   actRepo,
		statsRepo: statsRepo,
	}
}

// RecordProjectView records a project view
func (uc *AnalyticsUseCase) RecordProjectView(ctx context.Context, projectID, userID int64) error {
	view := entity.NewProjectView(projectID, userID)
	return uc.viewRepo.Record(ctx, view)
}

// GetProjectViews gets project views within a date range
func (uc *AnalyticsUseCase) GetProjectViews(ctx context.Context, projectID int64, startDate, endDate *time.Time) ([]*entity.ProjectView, int, error) {
	views, err := uc.viewRepo.GetByProjectID(ctx, projectID, startDate, endDate)
	if err != nil {
		return nil, 0, err
	}
	count, err := uc.viewRepo.CountByProjectID(ctx, projectID)
	if err != nil {
		return nil, 0, err
	}
	return views, count, nil
}

// RecordTaskActivity records a task activity
func (uc *AnalyticsUseCase) RecordTaskActivity(ctx context.Context, taskID, userID int64, action string) error {
	activity := entity.NewTaskActivity(taskID, userID, action)
	return uc.actRepo.Record(ctx, activity)
}

// GetTaskActivities gets activities for a task
func (uc *AnalyticsUseCase) GetTaskActivities(ctx context.Context, taskID int64) ([]*entity.TaskActivity, error) {
	return uc.actRepo.GetByTaskID(ctx, taskID)
}

// GetProjectStats gets stats for a project
func (uc *AnalyticsUseCase) GetProjectStats(ctx context.Context, projectID int64) (*entity.ProjectStats, error) {
	stats, err := uc.statsRepo.Get(ctx, projectID)
	if err != nil {
		return nil, ErrProjectStatsNotFound
	}
	return stats, nil
}

// UpdateProjectStats updates stats for a project
func (uc *AnalyticsUseCase) UpdateProjectStats(ctx context.Context, projectID int64, totalTasks int, completedTasks int) (*entity.ProjectStats, error) {
	stats := &entity.ProjectStats{
		ProjectID:      projectID,
		TotalTasks:     totalTasks,
		CompletedTasks: completedTasks,
	}
	stats.UpdateProgress()
	fmt.Println(stats)
	if err := uc.statsRepo.Upsert(ctx, stats); err != nil {
		return nil, err
	}
	return stats, nil
}

// GetDashboardStats gets dashboard statistics
func (uc *AnalyticsUseCase) GetDashboardStats(ctx context.Context) (*entity.DashboardStats, error) {
	allStats, err := uc.statsRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dashboard := &entity.DashboardStats{
		ProjectStats: allStats,
	}

	for _, stats := range allStats {
		dashboard.TotalProjects++
		if stats.ProgressPercent < 100 {
			dashboard.ActiveProjects++
		}
		dashboard.TotalTasks += stats.TotalTasks
		dashboard.CompletedTasks += stats.CompletedTasks
	}
	dashboard.PendingTasks = dashboard.TotalTasks - dashboard.CompletedTasks

	return dashboard, nil
}
