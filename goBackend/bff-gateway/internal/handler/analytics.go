package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/portfolio/proto/analytics"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AnalyticsHandler handles analytics endpoints
type AnalyticsHandler struct {
	analyticsClient pb.AnalyticsServiceClient
}

// NewAnalyticsHandler creates a new AnalyticsHandler
func NewAnalyticsHandler(conn *grpc.ClientConn) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsClient: pb.NewAnalyticsServiceClient(conn),
	}
}

func parseTimeOrNil(t string) *timestamppb.Timestamp {
	if t == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed)
}

// RecordProjectView records a project view
// POST /api/analytics/projects/:id/view
func (h *AnalyticsHandler) RecordProjectView(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID"})
		return
	}

	userIDVal, _ := c.Get("user_id")
	var userID int64
	if v, ok := userIDVal.(float64); ok {
		userID = int64(v)
	} else if v, ok := userIDVal.(int64); ok {
		userID = v
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.analyticsClient.RecordProjectView(ctx, &pb.RecordProjectViewRequest{
		ProjectId: projectID,
		UserId:    userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Record project view endpoint",
		"project_id": projectID,
		"user_id":    userID,
	})
}

// GetProjectViews returns project view statistics
// GET /api/analytics/projects/:id/views
func (h *AnalyticsHandler) GetProjectViews(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.analyticsClient.GetProjectViews(ctx, &pb.GetProjectViewsRequest{
		ProjectId: projectID,
		StartDate: parseTimeOrNil(startDate),
		EndDate:   parseTimeOrNil(endDate),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RecordTaskActivity records a task activity
// POST /api/analytics/tasks/:id/activity
func (h *AnalyticsHandler) RecordTaskActivity(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	userIDVal, _ := c.Get("user_id")
	var userID int64
	if v, ok := userIDVal.(float64); ok {
		userID = int64(v)
	} else if v, ok := userIDVal.(int64); ok {
		userID = v
	}

	var req struct {
		Action string `json:"action" binding:"required"` // created, updated, completed
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.analyticsClient.RecordTaskActivity(ctx, &pb.RecordTaskActivityRequest{
		TaskId: taskID,
		UserId: userID,
		Action: req.Action,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Record task activity endpoint",
		"task_id": taskID,
		"user_id": userID,
		"action":  req.Action,
	})
}

// GetTaskActivities returns task activity log
// GET /api/analytics/tasks/:id/activities
func (h *AnalyticsHandler) GetTaskActivities(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.analyticsClient.GetTaskActivities(ctx, &pb.GetTaskActivitiesRequest{
		TaskId: taskID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Activities)
}

// GetProjectStats returns project statistics
// GET /api/analytics/projects/:id/stats
func (h *AnalyticsHandler) GetProjectStats(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.analyticsClient.GetProjectStats(ctx, &pb.GetProjectStatsRequest{
		ProjectId: projectID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Stats)
}

// GetDashboardStats returns dashboard statistics
// GET /api/analytics/dashboard
func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	var userID int64
	if exists {
		if v, ok := userIDVal.(float64); ok {
			userID = int64(v)
		} else if v, ok := userIDVal.(int64); ok {
			userID = v
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.analyticsClient.GetDashboardStats(ctx, &pb.GetDashboardStatsRequest{
		UserId: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
