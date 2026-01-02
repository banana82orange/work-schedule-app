package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/portfolio/task-service/internal/domain/entity"
	"github.com/portfolio/task-service/internal/domain/repository"
)

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrSubtaskNotFound = errors.New("subtask not found")
	ErrCommentNotFound = errors.New("comment not found")
)

// TaskUseCase handles task business logic
type TaskUseCase struct {
	taskRepo       repository.TaskRepository
	subtaskRepo    repository.SubtaskRepository
	commentRepo    repository.CommentRepository
	attachmentRepo repository.AttachmentRepository
	tagRepo        repository.TagRepository
	taskTagRepo    repository.TaskTagRepository
}

// NewTaskUseCase creates a new TaskUseCase
func NewTaskUseCase(
	taskRepo repository.TaskRepository,
	subtaskRepo repository.SubtaskRepository,
	commentRepo repository.CommentRepository,
	attachmentRepo repository.AttachmentRepository,
	tagRepo repository.TagRepository,
	taskTagRepo repository.TaskTagRepository,
) *TaskUseCase {
	return &TaskUseCase{
		taskRepo:       taskRepo,
		subtaskRepo:    subtaskRepo,
		commentRepo:    commentRepo,
		attachmentRepo: attachmentRepo,
		tagRepo:        tagRepo,
		taskTagRepo:    taskTagRepo,
	}
}

// CreateTask creates a new task
func (uc *TaskUseCase) CreateTask(ctx context.Context, projectID int64, title, description, status string, priority int, assignedTo int64, dueDate *time.Time) (*entity.Task, error) {
	fmt.Println("CreateTask")
	fmt.Println(projectID, title, description, status, priority, assignedTo, dueDate)
	task := entity.NewTask(projectID, title, description, status, priority, assignedTo, dueDate)
	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

// GetTask retrieves a task by ID with all related data
func (uc *TaskUseCase) GetTask(ctx context.Context, id int64) (*entity.Task, error) {
	task, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Load subtasks
	subtasks, _ := uc.subtaskRepo.GetByTaskID(ctx, id)
	task.Subtasks = subtasks

	// Load tags
	tags, _ := uc.taskTagRepo.GetByTaskID(ctx, id)
	task.Tags = tags

	return task, nil
}

// UpdateTask updates a task
func (uc *TaskUseCase) UpdateTask(ctx context.Context, id int64, title, description, status string, priority int, assignedTo int64, dueDate *time.Time) (*entity.Task, error) {
	task, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if status != "" {
		task.Status = status
	}
	if priority > 0 {
		task.Priority = priority
	}
	if assignedTo > 0 {
		task.AssignedTo = &assignedTo
	}
	if dueDate != nil {
		task.DueDate = dueDate
	}
	task.UpdatedAt = time.Now()

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return uc.GetTask(ctx, id)
}

// DeleteTask deletes a task
func (uc *TaskUseCase) DeleteTask(ctx context.Context, id int64) error {
	return uc.taskRepo.Delete(ctx, id)
}

// ListTasks lists tasks with filters
func (uc *TaskUseCase) ListTasks(ctx context.Context, projectID int64, page, limit int, status string, assignedTo int64) ([]*entity.Task, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.taskRepo.List(ctx, projectID, page, limit, status, assignedTo)
}

// SubtaskUseCase handles subtask business logic
type SubtaskUseCase struct {
	subtaskRepo repository.SubtaskRepository
}

// NewSubtaskUseCase creates a new SubtaskUseCase
func NewSubtaskUseCase(subtaskRepo repository.SubtaskRepository) *SubtaskUseCase {
	return &SubtaskUseCase{subtaskRepo: subtaskRepo}
}

// CreateSubtask creates a new subtask
func (uc *SubtaskUseCase) CreateSubtask(ctx context.Context, taskID int64, title string, assignedTo int64, dueDate *time.Time) (*entity.Subtask, error) {
	subtask := entity.NewSubtask(taskID, title, assignedTo, dueDate)
	if err := uc.subtaskRepo.Create(ctx, subtask); err != nil {
		return nil, err
	}
	return subtask, nil
}

// UpdateSubtask updates a subtask
func (uc *SubtaskUseCase) UpdateSubtask(ctx context.Context, id int64, title, status string, assignedTo int64, dueDate *time.Time) (*entity.Subtask, error) {
	subtask, err := uc.subtaskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSubtaskNotFound
	}

	if title != "" {
		subtask.Title = title
	}
	if status != "" {
		subtask.Status = status
	}
	if assignedTo > 0 {
		subtask.AssignedTo = assignedTo
	}
	if dueDate != nil {
		subtask.DueDate = dueDate
	}
	subtask.UpdatedAt = time.Now()

	if err := uc.subtaskRepo.Update(ctx, subtask); err != nil {
		return nil, err
	}

	return subtask, nil
}

// DeleteSubtask deletes a subtask
func (uc *SubtaskUseCase) DeleteSubtask(ctx context.Context, id int64) error {
	return uc.subtaskRepo.Delete(ctx, id)
}

// GetSubtasks gets all subtasks for a task
func (uc *SubtaskUseCase) GetSubtasks(ctx context.Context, taskID int64) ([]*entity.Subtask, error) {
	return uc.subtaskRepo.GetByTaskID(ctx, taskID)
}

// CommentUseCase handles comment business logic
type CommentUseCase struct {
	commentRepo repository.CommentRepository
}

// NewCommentUseCase creates a new CommentUseCase
func NewCommentUseCase(commentRepo repository.CommentRepository) *CommentUseCase {
	return &CommentUseCase{commentRepo: commentRepo}
}

// AddComment adds a comment to a task
func (uc *CommentUseCase) AddComment(ctx context.Context, taskID, userID int64, comment string) (*entity.TaskComment, error) {
	taskComment := entity.NewTaskComment(taskID, userID, comment)
	if err := uc.commentRepo.Create(ctx, taskComment); err != nil {
		return nil, err
	}
	return taskComment, nil
}

// DeleteComment deletes a comment
func (uc *CommentUseCase) DeleteComment(ctx context.Context, id int64) error {
	return uc.commentRepo.Delete(ctx, id)
}

// GetComments gets all comments for a task
func (uc *CommentUseCase) GetComments(ctx context.Context, taskID int64) ([]*entity.TaskComment, error) {
	return uc.commentRepo.GetByTaskID(ctx, taskID)
}

// AttachmentUseCase handles attachment business logic
type AttachmentUseCase struct {
	attachmentRepo repository.AttachmentRepository
}

// NewAttachmentUseCase creates a new AttachmentUseCase
func NewAttachmentUseCase(attachmentRepo repository.AttachmentRepository) *AttachmentUseCase {
	return &AttachmentUseCase{attachmentRepo: attachmentRepo}
}

// AddAttachment adds an attachment to a task
func (uc *AttachmentUseCase) AddAttachment(ctx context.Context, taskID int64, fileURL string) (*entity.TaskAttachment, error) {
	attachment := entity.NewTaskAttachment(taskID, fileURL)
	if err := uc.attachmentRepo.Create(ctx, attachment); err != nil {
		return nil, err
	}
	return attachment, nil
}

// DeleteAttachment deletes an attachment
func (uc *AttachmentUseCase) DeleteAttachment(ctx context.Context, id int64) error {
	return uc.attachmentRepo.Delete(ctx, id)
}

// GetAttachments gets all attachments for a task
func (uc *AttachmentUseCase) GetAttachments(ctx context.Context, taskID int64) ([]*entity.TaskAttachment, error) {
	return uc.attachmentRepo.GetByTaskID(ctx, taskID)
}

// TagUseCase handles tag business logic
type TagUseCase struct {
	tagRepo     repository.TagRepository
	taskTagRepo repository.TaskTagRepository
}

// NewTagUseCase creates a new TagUseCase
func NewTagUseCase(tagRepo repository.TagRepository, taskTagRepo repository.TaskTagRepository) *TagUseCase {
	return &TagUseCase{
		tagRepo:     tagRepo,
		taskTagRepo: taskTagRepo,
	}
}

// CreateTag creates a new tag
func (uc *TagUseCase) CreateTag(ctx context.Context, name string) (*entity.TaskTag, error) {
	tag := &entity.TaskTag{Name: name}
	if err := uc.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

// ListTags lists all tags
func (uc *TagUseCase) ListTags(ctx context.Context) ([]*entity.TaskTag, error) {
	return uc.tagRepo.List(ctx)
}

// AddTaskTag adds a tag to a task
func (uc *TagUseCase) AddTaskTag(ctx context.Context, taskID, tagID int64) error {
	return uc.taskTagRepo.Add(ctx, taskID, tagID)
}

// RemoveTaskTag removes a tag from a task
func (uc *TagUseCase) RemoveTaskTag(ctx context.Context, taskID, tagID int64) error {
	return uc.taskTagRepo.Remove(ctx, taskID, tagID)
}
