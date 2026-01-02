package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/portfolio/task-service/internal/domain/entity"
)

// PostgresTaskRepository implements TaskRepository
type PostgresTaskRepository struct {
	db *sql.DB
}

// NewPostgresTaskRepository creates a new PostgresTaskRepository
func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

// Create creates a new task
func (r *PostgresTaskRepository) Create(ctx context.Context, task *entity.Task) error {
	query := `
		INSERT INTO tasks (project_id, title, description, status, priority, assigned_to, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, DATE($7), $8, $9)
		RETURNING id
	`
	return r.db.QueryRowContext(
		ctx, query,
		task.ProjectID, task.Title, task.Description, task.Status,
		task.Priority, task.AssignedTo, task.DueDate, task.CreatedAt, task.UpdatedAt,
	).Scan(&task.ID)
}

// GetByID gets a task by ID
func (r *PostgresTaskRepository) GetByID(ctx context.Context, id int64) (*entity.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, priority, assigned_to, due_date, created_at, updated_at
		FROM tasks WHERE id = $1
	`
	var description sql.NullString
	task := &entity.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.ProjectID, &task.Title, &description,
		&task.Status, &task.Priority, &task.AssignedTo, &task.DueDate,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if description.Valid {
		task.Description = description.String
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Update updates a task
func (r *PostgresTaskRepository) Update(ctx context.Context, task *entity.Task) error {
	query := `
		UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4,
		assigned_to = $5, due_date = $6, updated_at = $7 WHERE id = $8
	`
	task.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Status, task.Priority,
		task.AssignedTo, task.DueDate, task.UpdatedAt, task.ID,
	)
	return err
}

// Delete deletes a task
func (r *PostgresTaskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List lists tasks with filters
func (r *PostgresTaskRepository) List(ctx context.Context, projectID int64, page, limit int, status string, assignedTo int64) ([]*entity.Task, int, error) {
	offset := (page - 1) * limit

	// Build dynamic query
	baseQuery := `FROM tasks WHERE project_id = $1`
	args := []interface{}{projectID}
	argIndex := 2

	if status != "" {
		baseQuery += ` AND status = $` + string(rune('0'+argIndex))
		args = append(args, status)
		argIndex++
	}
	if assignedTo > 0 {
		baseQuery += ` AND assigned_to = $` + string(rune('0'+argIndex))
		args = append(args, assignedTo)
		argIndex++
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) ` + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get tasks
	selectQuery := `SELECT id, project_id, title, description, status, priority, assigned_to, due_date, created_at, updated_at ` + baseQuery + ` ORDER BY priority, due_date LIMIT $` + string(rune('0'+argIndex)) + ` OFFSET $` + string(rune('0'+argIndex+1))
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*entity.Task
	for rows.Next() {
		task := &entity.Task{}
		var description sql.NullString
		if err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Title, &description,
			&task.Status, &task.Priority, &task.AssignedTo, &task.DueDate,
			&task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		if description.Valid {
			task.Description = description.String
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

// PostgresSubtaskRepository implements SubtaskRepository
type PostgresSubtaskRepository struct {
	db *sql.DB
}

// NewPostgresSubtaskRepository creates a new repository
func NewPostgresSubtaskRepository(db *sql.DB) *PostgresSubtaskRepository {
	return &PostgresSubtaskRepository{db: db}
}

// Create creates a new subtask
func (r *PostgresSubtaskRepository) Create(ctx context.Context, subtask *entity.Subtask) error {
	query := `
		INSERT INTO subtasks (task_id, title, status, assigned_to, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		subtask.TaskID, subtask.Title, subtask.Status, subtask.AssignedTo,
		subtask.DueDate, subtask.CreatedAt, subtask.UpdatedAt,
	).Scan(&subtask.ID)
}

// GetByID gets a subtask by ID
func (r *PostgresSubtaskRepository) GetByID(ctx context.Context, id int64) (*entity.Subtask, error) {
	query := `SELECT id, task_id, title, status, assigned_to, due_date, created_at, updated_at FROM subtasks WHERE id = $1`
	subtask := &entity.Subtask{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subtask.ID, &subtask.TaskID, &subtask.Title, &subtask.Status,
		&subtask.AssignedTo, &subtask.DueDate, &subtask.CreatedAt, &subtask.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return subtask, nil
}

// Update updates a subtask
func (r *PostgresSubtaskRepository) Update(ctx context.Context, subtask *entity.Subtask) error {
	query := `UPDATE subtasks SET title = $1, status = $2, assigned_to = $3, due_date = $4, updated_at = $5 WHERE id = $6`
	subtask.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, subtask.Title, subtask.Status, subtask.AssignedTo, subtask.DueDate, subtask.UpdatedAt, subtask.ID)
	return err
}

// Delete deletes a subtask
func (r *PostgresSubtaskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM subtasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByTaskID gets all subtasks for a task
func (r *PostgresSubtaskRepository) GetByTaskID(ctx context.Context, taskID int64) ([]*entity.Subtask, error) {
	query := `SELECT id, task_id, title, status, assigned_to, due_date, created_at, updated_at FROM subtasks WHERE task_id = $1`
	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtasks []*entity.Subtask
	for rows.Next() {
		subtask := &entity.Subtask{}
		if err := rows.Scan(&subtask.ID, &subtask.TaskID, &subtask.Title, &subtask.Status, &subtask.AssignedTo, &subtask.DueDate, &subtask.CreatedAt, &subtask.UpdatedAt); err != nil {
			return nil, err
		}
		subtasks = append(subtasks, subtask)
	}
	return subtasks, nil
}

// PostgresCommentRepository implements CommentRepository
type PostgresCommentRepository struct {
	db *sql.DB
}

// NewPostgresCommentRepository creates a new repository
func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

// Create creates a new comment
func (r *PostgresCommentRepository) Create(ctx context.Context, comment *entity.TaskComment) error {
	query := `INSERT INTO task_comments (task_id, user_id, comment, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, comment.TaskID, comment.UserID, comment.Comment, comment.CreatedAt).Scan(&comment.ID)
}

// GetByID gets a comment by ID
func (r *PostgresCommentRepository) GetByID(ctx context.Context, id int64) (*entity.TaskComment, error) {
	query := `SELECT id, task_id, user_id, comment, created_at FROM task_comments WHERE id = $1`
	comment := &entity.TaskComment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.Comment, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// Delete deletes a comment
func (r *PostgresCommentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM task_comments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByTaskID gets all comments for a task
func (r *PostgresCommentRepository) GetByTaskID(ctx context.Context, taskID int64) ([]*entity.TaskComment, error) {
	query := `SELECT id, task_id, user_id, comment, created_at FROM task_comments WHERE task_id = $1 ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.TaskComment
	for rows.Next() {
		comment := &entity.TaskComment{}
		if err := rows.Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.Comment, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// PostgresAttachmentRepository implements AttachmentRepository
type PostgresAttachmentRepository struct {
	db *sql.DB
}

// NewPostgresAttachmentRepository creates a new repository
func NewPostgresAttachmentRepository(db *sql.DB) *PostgresAttachmentRepository {
	return &PostgresAttachmentRepository{db: db}
}

// Create creates a new attachment
func (r *PostgresAttachmentRepository) Create(ctx context.Context, attachment *entity.TaskAttachment) error {
	query := `INSERT INTO task_attachments (task_id, file_url, uploaded_at) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, attachment.TaskID, attachment.FileURL, attachment.UploadedAt).Scan(&attachment.ID)
}

// GetByID gets an attachment by ID
func (r *PostgresAttachmentRepository) GetByID(ctx context.Context, id int64) (*entity.TaskAttachment, error) {
	query := `SELECT id, task_id, file_url, uploaded_at FROM task_attachments WHERE id = $1`
	attachment := &entity.TaskAttachment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&attachment.ID, &attachment.TaskID, &attachment.FileURL, &attachment.UploadedAt)
	if err != nil {
		return nil, err
	}
	return attachment, nil
}

// Delete deletes an attachment
func (r *PostgresAttachmentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM task_attachments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByTaskID gets all attachments for a task
func (r *PostgresAttachmentRepository) GetByTaskID(ctx context.Context, taskID int64) ([]*entity.TaskAttachment, error) {
	query := `SELECT id, task_id, file_url, uploaded_at FROM task_attachments WHERE task_id = $1`
	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []*entity.TaskAttachment
	for rows.Next() {
		attachment := &entity.TaskAttachment{}
		if err := rows.Scan(&attachment.ID, &attachment.TaskID, &attachment.FileURL, &attachment.UploadedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}

// PostgresTagRepository implements TagRepository
type PostgresTagRepository struct {
	db *sql.DB
}

// NewPostgresTagRepository creates a new repository
func NewPostgresTagRepository(db *sql.DB) *PostgresTagRepository {
	return &PostgresTagRepository{db: db}
}

// Create creates a new tag
func (r *PostgresTagRepository) Create(ctx context.Context, tag *entity.TaskTag) error {
	query := `INSERT INTO task_tags (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, tag.Name).Scan(&tag.ID)
}

// GetByID gets a tag by ID
func (r *PostgresTagRepository) GetByID(ctx context.Context, id int64) (*entity.TaskTag, error) {
	query := `SELECT id, name FROM task_tags WHERE id = $1`
	tag := &entity.TaskTag{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&tag.ID, &tag.Name)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// List lists all tags
func (r *PostgresTagRepository) List(ctx context.Context) ([]*entity.TaskTag, error) {
	query := `SELECT id, name FROM task_tags ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entity.TaskTag
	for rows.Next() {
		tag := &entity.TaskTag{}
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// PostgresTaskTagRepository implements TaskTagRepository
type PostgresTaskTagRepository struct {
	db *sql.DB
}

// NewPostgresTaskTagRepository creates a new repository
func NewPostgresTaskTagRepository(db *sql.DB) *PostgresTaskTagRepository {
	return &PostgresTaskTagRepository{db: db}
}

// Add adds a tag to a task
func (r *PostgresTaskTagRepository) Add(ctx context.Context, taskID, tagID int64) error {
	query := `INSERT INTO task_tag_mapping (task_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, taskID, tagID)
	return err
}

// Remove removes a tag from a task
func (r *PostgresTaskTagRepository) Remove(ctx context.Context, taskID, tagID int64) error {
	query := `DELETE FROM task_tag_mapping WHERE task_id = $1 AND tag_id = $2`
	_, err := r.db.ExecContext(ctx, query, taskID, tagID)
	return err
}

// GetByTaskID gets all tags for a task
func (r *PostgresTaskTagRepository) GetByTaskID(ctx context.Context, taskID int64) ([]*entity.TaskTag, error) {
	query := `SELECT t.id, t.name FROM task_tags t INNER JOIN task_tag_mapping m ON t.id = m.tag_id WHERE m.task_id = $1`
	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entity.TaskTag
	for rows.Next() {
		tag := &entity.TaskTag{}
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
