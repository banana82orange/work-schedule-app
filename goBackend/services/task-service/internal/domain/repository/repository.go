package repository

import (
	"context"

	"github.com/portfolio/task-service/internal/domain/entity"
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id int64) (*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, projectID int64, page, limit int, status string, assignedTo int64) ([]*entity.Task, int, error)
}

// SubtaskRepository defines the interface for subtask data access
type SubtaskRepository interface {
	Create(ctx context.Context, subtask *entity.Subtask) error
	GetByID(ctx context.Context, id int64) (*entity.Subtask, error)
	Update(ctx context.Context, subtask *entity.Subtask) error
	Delete(ctx context.Context, id int64) error
	GetByTaskID(ctx context.Context, taskID int64) ([]*entity.Subtask, error)
}

// CommentRepository defines the interface for comment data access
type CommentRepository interface {
	Create(ctx context.Context, comment *entity.TaskComment) error
	GetByID(ctx context.Context, id int64) (*entity.TaskComment, error)
	Delete(ctx context.Context, id int64) error
	GetByTaskID(ctx context.Context, taskID int64) ([]*entity.TaskComment, error)
}

// AttachmentRepository defines the interface for attachment data access
type AttachmentRepository interface {
	Create(ctx context.Context, attachment *entity.TaskAttachment) error
	GetByID(ctx context.Context, id int64) (*entity.TaskAttachment, error)
	Delete(ctx context.Context, id int64) error
	GetByTaskID(ctx context.Context, taskID int64) ([]*entity.TaskAttachment, error)
}

// TagRepository defines the interface for tag data access
type TagRepository interface {
	Create(ctx context.Context, tag *entity.TaskTag) error
	GetByID(ctx context.Context, id int64) (*entity.TaskTag, error)
	List(ctx context.Context) ([]*entity.TaskTag, error)
}

// TaskTagRepository defines the interface for task-tag relationship
type TaskTagRepository interface {
	Add(ctx context.Context, taskID, tagID int64) error
	Remove(ctx context.Context, taskID, tagID int64) error
	GetByTaskID(ctx context.Context, taskID int64) ([]*entity.TaskTag, error)
}
