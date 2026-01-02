package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/portfolio/proto/task"
	"google.golang.org/grpc"
)

// TaskHandler handles task endpoints
type TaskHandler struct {
	taskClient pb.TaskServiceClient
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(conn *grpc.ClientConn) *TaskHandler {
	return &TaskHandler{
		taskClient: pb.NewTaskServiceClient(conn),
	}
}

// CreateTaskRequest represents create task request
type CreateTaskRequest struct {
	ProjectID   int64  `json:"project_id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    int32  `json:"priority"`
	AssignedTo  int64  `json:"assigned_to"`
	DueDate     string `json:"due_date"`
}


// CreateTask creates a new task
// POST /api/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.CreateTask(ctx, &pb.CreateTaskRequest{
		ProjectId:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssignedTo:  req.AssignedTo,
		DueDate:     parseTime(req.DueDate),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Task)
}

// GetTask returns a task by ID
// GET /api/tasks/:id
func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.GetTask(ctx, &pb.GetTaskRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Task)
}

// UpdateTask updates a task
// PUT /api/tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssignedTo:  req.AssignedTo,
		DueDate:     parseTime(req.DueDate),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Task)
}

// DeleteTask deletes a task
// DELETE /api/tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.taskClient.DeleteTask(ctx, &pb.DeleteTaskRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// ListTasks returns list of tasks
// GET /api/tasks
func (h *TaskHandler) ListTasks(c *gin.Context) {
	// page := c.DefaultQuery("page", "1")
	// limit := c.DefaultQuery("limit", "10")
	status := c.Query("status")
	projectIDStr := c.Query("project_id")
	var projectID int64
	if projectIDStr != "" {
		projectID, _ = strconv.ParseInt(projectIDStr, 10, 64)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.ListTasks(ctx, &pb.ListTasksRequest{
		ProjectId: projectID,
		Page:      1,
		Limit:     100, // fetching more for now
		Status:    status,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Tasks)
}

// CreateSubtask creates a new subtask
// POST /api/tasks/:id/subtasks
func (h *TaskHandler) CreateSubtask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	var req struct {
		Title      string `json:"title" binding:"required"`
		AssignedTo int64  `json:"assigned_to"`
		DueDate    string `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.CreateSubtask(ctx, &pb.CreateSubtaskRequest{
		TaskId:     taskID,
		Title:      req.Title,
		AssignedTo: req.AssignedTo,
		DueDate:    parseTime(req.DueDate),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Subtask)
}

// ListSubtasks returns list of subtasks
// GET /api/tasks/:id/subtasks
func (h *TaskHandler) ListSubtasks(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.ListSubtasks(ctx, &pb.ListSubtasksRequest{TaskId: taskID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Subtasks)
}

// AddComment adds a comment to task
// POST /api/tasks/:id/comments
func (h *TaskHandler) AddComment(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	var req struct {
		Comment string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id") // Assuming set by middleware
	var uid int64
	// Handle both float64 (json default) and int
	if v, ok := userID.(float64); ok {
		uid = int64(v)
	} else if v, ok := userID.(int); ok {
		uid = int64(v)
	} else if v, ok := userID.(int64); ok {
		uid = v
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.AddComment(ctx, &pb.AddCommentRequest{
		TaskId:  taskID,
		UserId:  uid,
		Comment: req.Comment,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Comment)
}

// ListComments returns list of comments
// GET /api/tasks/:id/comments
func (h *TaskHandler) ListComments(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.ListComments(ctx, &pb.ListCommentsRequest{TaskId: taskID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Comments)
}

// AddAttachment adds attachment to task
// POST /api/tasks/:id/attachments
func (h *TaskHandler) AddAttachment(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	var req struct {
		FileURL string `json:"file_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.AddAttachment(ctx, &pb.AddAttachmentRequest{
		TaskId:  taskID,
		FileUrl: req.FileURL,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Attachment)
}

// ListAttachments returns list of attachments
// GET /api/tasks/:id/attachments
func (h *TaskHandler) ListAttachments(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.ListAttachments(ctx, &pb.ListAttachmentsRequest{TaskId: taskID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Attachments)
}

// AddTag adds a tag to task
// POST /api/tasks/:id/tags
func (h *TaskHandler) AddTaskTag(c *gin.Context) {
	// Just alias or implement if different
	h.AddTag(c)
}

// CreateTag creates a new tag
// POST /api/tags
func (h *TaskHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.CreateTag(ctx, &pb.CreateTagRequest{Name: req.Name})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Tag)
}

// ListTags returns all tags
// GET /api/tags
func (h *TaskHandler) ListTags(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.taskClient.ListTags(ctx, &pb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Tags)
}

// AddTag implementation
func (h *TaskHandler) AddTag(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	var req struct {
		TagID int64 `json:"tag_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.taskClient.AddTaskTag(ctx, &pb.AddTaskTagRequest{
		TaskId: taskID,
		TagId:  req.TagID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added to task"})
}
