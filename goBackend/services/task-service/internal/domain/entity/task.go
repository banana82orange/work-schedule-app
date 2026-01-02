package entity

import (
	"time"
)

// Task represents a task entity
type Task struct {
	ID          int64       `json:"id"`
	ProjectID   int64       `json:"project_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      string      `json:"status"` // Todo, InProgress, Done
	Priority    int         `json:"priority"`
	AssignedTo  *int64      `json:"assigned_to,omitempty"`
	DueDate     *time.Time  `json:"due_date,omitempty"`
	Subtasks    []*Subtask  `json:"subtasks,omitempty"`
	Tags        []*TaskTag  `json:"tags,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// NewTask creates a new task entity
func NewTask(projectID int64, title, description, status string, priority int, assignedTo int64, dueDate *time.Time) *Task {
	now := time.Now()
	if status == "" {
		status = StatusTodo
	}
	if priority == 0 {
		priority = 3
	}

	var assignedToPtr *int64
	if assignedTo != 0 {
		assignedToPtr = &assignedTo
	}

	return &Task{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		AssignedTo:  assignedToPtr,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Task statuses
const (
	StatusTodo       = "Todo"
	StatusInProgress = "InProgress"
	StatusDone       = "Done"
)

// ValidTaskStatuses returns all valid task statuses
func ValidTaskStatuses() []string {
	return []string{StatusTodo, StatusInProgress, StatusDone}
}

// Subtask represents a subtask entity
type Subtask struct {
	ID         int64      `json:"id"`
	TaskID     int64      `json:"task_id"`
	Title      string     `json:"title"`
	Status     string     `json:"status"`
	AssignedTo int64      `json:"assigned_to"`
	DueDate    *time.Time `json:"due_date,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// NewSubtask creates a new subtask entity
func NewSubtask(taskID int64, title string, assignedTo int64, dueDate *time.Time) *Subtask {
	now := time.Now()
	return &Subtask{
		TaskID:     taskID,
		Title:      title,
		Status:     StatusTodo,
		AssignedTo: assignedTo,
		DueDate:    dueDate,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// TaskComment represents a task comment
type TaskComment struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	UserID    int64     `json:"user_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTaskComment creates a new task comment
func NewTaskComment(taskID, userID int64, comment string) *TaskComment {
	return &TaskComment{
		TaskID:    taskID,
		UserID:    userID,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
}

// TaskAttachment represents a task attachment
type TaskAttachment struct {
	ID         int64     `json:"id"`
	TaskID     int64     `json:"task_id"`
	FileURL    string    `json:"file_url"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// NewTaskAttachment creates a new task attachment
func NewTaskAttachment(taskID int64, fileURL string) *TaskAttachment {
	return &TaskAttachment{
		TaskID:     taskID,
		FileURL:    fileURL,
		UploadedAt: time.Now(),
	}
}

// TaskTag represents a task tag
type TaskTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TaskTagMapping represents task-tag relationship
type TaskTagMapping struct {
	TaskID int64 `json:"task_id"`
	TagID  int64 `json:"tag_id"`
}
