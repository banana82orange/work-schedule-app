package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/portfolio/analytics-service/internal/domain/entity"
)

// PostgresProjectViewRepository implements ProjectViewRepository
type PostgresProjectViewRepository struct {
	db *sql.DB
}

// NewPostgresProjectViewRepository creates a new repository
func NewPostgresProjectViewRepository(db *sql.DB) *PostgresProjectViewRepository {
	return &PostgresProjectViewRepository{db: db}
}

// Record records a project view
func (r *PostgresProjectViewRepository) Record(ctx context.Context, view *entity.ProjectView) error {
	query := `INSERT INTO project_views (project_id, user_id, viewed_at) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, view.ProjectID, view.UserID, view.ViewedAt).Scan(&view.ID)
}

// GetByProjectID gets project views with optional date range
func (r *PostgresProjectViewRepository) GetByProjectID(ctx context.Context, projectID int64, startDate, endDate *time.Time) ([]*entity.ProjectView, error) {
	query := `SELECT id, project_id, user_id, viewed_at FROM project_views WHERE project_id = $1`
	args := []interface{}{projectID}
	argIndex := 2

	if startDate != nil {
		query += ` AND viewed_at >= $` + string(rune('0'+argIndex))
		args = append(args, startDate)
		argIndex++
	}
	if endDate != nil {
		query += ` AND viewed_at <= $` + string(rune('0'+argIndex))
		args = append(args, endDate)
	}
	query += ` ORDER BY viewed_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []*entity.ProjectView
	for rows.Next() {
		view := &entity.ProjectView{}
		if err := rows.Scan(&view.ID, &view.ProjectID, &view.UserID, &view.ViewedAt); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

// CountByProjectID counts total views for a project
func (r *PostgresProjectViewRepository) CountByProjectID(ctx context.Context, projectID int64) (int, error) {
	query := `SELECT COUNT(*) FROM project_views WHERE project_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(&count)
	return count, err
}

// PostgresTaskActivityRepository implements TaskActivityRepository
type PostgresTaskActivityRepository struct {
	db *sql.DB
}

// NewPostgresTaskActivityRepository creates a new repository
func NewPostgresTaskActivityRepository(db *sql.DB) *PostgresTaskActivityRepository {
	return &PostgresTaskActivityRepository{db: db}
}

// Record records a task activity
func (r *PostgresTaskActivityRepository) Record(ctx context.Context, activity *entity.TaskActivity) error {
	query := `INSERT INTO task_activity (task_id, user_id, action, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, activity.TaskID, activity.UserID, activity.Action, activity.CreatedAt).Scan(&activity.ID)
}

// GetByTaskID gets activities for a task
func (r *PostgresTaskActivityRepository) GetByTaskID(ctx context.Context, taskID int64) ([]*entity.TaskActivity, error) {
	query := `SELECT id, task_id, user_id, action, created_at FROM task_activity WHERE task_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*entity.TaskActivity
	for rows.Next() {
		activity := &entity.TaskActivity{}
		if err := rows.Scan(&activity.ID, &activity.TaskID, &activity.UserID, &activity.Action, &activity.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

// GetByProjectID gets activities for all tasks in a project
func (r *PostgresTaskActivityRepository) GetByProjectID(ctx context.Context, projectID int64) ([]*entity.TaskActivity, error) {
	query := `
		SELECT ta.id, ta.task_id, ta.user_id, ta.action, ta.created_at
		FROM task_activity ta
		INNER JOIN tasks t ON ta.task_id = t.id
		WHERE t.project_id = $1
		ORDER BY ta.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*entity.TaskActivity
	for rows.Next() {
		activity := &entity.TaskActivity{}
		if err := rows.Scan(&activity.ID, &activity.TaskID, &activity.UserID, &activity.Action, &activity.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

// PostgresProjectStatsRepository implements ProjectStatsRepository
type PostgresProjectStatsRepository struct {
	db *sql.DB
}

// NewPostgresProjectStatsRepository creates a new repository
func NewPostgresProjectStatsRepository(db *sql.DB) *PostgresProjectStatsRepository {
	return &PostgresProjectStatsRepository{db: db}
}

// Get gets stats for a project
func (r *PostgresProjectStatsRepository) Get(ctx context.Context, projectID int64) (*entity.ProjectStats, error) {
	query := `SELECT project_id, total_tasks, completed_tasks, progress_percent, last_updated FROM project_stats WHERE project_id = $1`
	stats := &entity.ProjectStats{}
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&stats.ProjectID, &stats.TotalTasks, &stats.CompletedTasks,
		&stats.ProgressPercent, &stats.LastUpdated,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// Upsert inserts or updates project stats
func (r *PostgresProjectStatsRepository) Upsert(ctx context.Context, stats *entity.ProjectStats) error {
	query := `
		INSERT INTO project_stats (project_id, total_tasks, completed_tasks, progress_percent, last_updated)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (project_id) DO UPDATE SET
			total_tasks = $2, completed_tasks = $3, progress_percent = $4, last_updated = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		stats.ProjectID, stats.TotalTasks, stats.CompletedTasks,
		stats.ProgressPercent, time.Now(),
	)
	return err
}

// GetAll gets all project stats
func (r *PostgresProjectStatsRepository) GetAll(ctx context.Context) ([]*entity.ProjectStats, error) {
	query := `SELECT project_id, total_tasks, completed_tasks, progress_percent, last_updated FROM project_stats`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allStats []*entity.ProjectStats
	for rows.Next() {
		stats := &entity.ProjectStats{}
		if err := rows.Scan(&stats.ProjectID, &stats.TotalTasks, &stats.CompletedTasks, &stats.ProgressPercent, &stats.LastUpdated); err != nil {
			return nil, err
		}
		allStats = append(allStats, stats)
	}
	return allStats, nil
}
