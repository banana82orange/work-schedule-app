package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/portfolio/proto/project"
	"google.golang.org/grpc"
)

// ProjectHandler handles project endpoints
type ProjectHandler struct {
	projectClient pb.ProjectServiceClient
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(conn *grpc.ClientConn) *ProjectHandler {
	return &ProjectHandler{
		projectClient: pb.NewProjectServiceClient(conn),
	}
}

// CreateProjectRequest represents create project request
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Status      string `json:"status"`
}



// CreateProject creates a new project
// POST /api/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.projectClient.CreateProject(ctx, &pb.CreateProjectRequest{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   parseTime(req.StartDate),
		EndDate:     parseTime(req.EndDate),
		Status:      req.Status,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Project)
}

// GetProject returns a project by ID
// GET /api/projects/:id
func (h *ProjectHandler) GetProject(c *gin.Context) {
	var req struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.projectClient.GetProject(ctx, &pb.GetProjectRequest{Id: req.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Project)
}

// UpdateProject updates a project
// PUT /api/projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStruct := struct {
		ID int64 `uri:"id" binding:"required"`
	}{}
	if err := c.ShouldBindUri(&idStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.projectClient.UpdateProject(ctx, &pb.UpdateProjectRequest{
		Id:          idStruct.ID,
		Name:        req.Name,
		Description: req.Description,
		StartDate:   parseTime(req.StartDate),
		EndDate:     parseTime(req.EndDate),
		Status:      req.Status,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Project)
}

// DeleteProject deletes a project
// DELETE /api/projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	var req struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.projectClient.DeleteProject(ctx, &pb.DeleteProjectRequest{Id: req.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// ListProjects returns list of projects
// GET /api/projects
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	// page := c.DefaultQuery("page", "1")
	// limit := c.DefaultQuery("limit", "10")
	status := c.Query("status")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.projectClient.ListProjects(ctx, &pb.ListProjectsRequest{
		Page:   1, // Simplification
		Limit:  10,
		Status: status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Projects)
}

// AddSkill adds a skill to project
// POST /api/projects/:id/skills
func (h *ProjectHandler) AddSkill(c *gin.Context) {
	var uri struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		SkillID int64 `json:"skill_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.projectClient.AddProjectSkill(ctx, &pb.AddProjectSkillRequest{
		ProjectId: uri.ID,
		SkillId:   req.SkillID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill added to project"})
}

// AddTech adds technology to project
// POST /api/projects/:id/tech
func (h *ProjectHandler) AddTech(c *gin.Context) {
	var uri struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		TechName string `json:"tech_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.projectClient.AddProjectTech(ctx, &pb.AddProjectTechRequest{
		ProjectId: uri.ID,
		TechName:  req.TechName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tech added to project"})
}

// AddImage adds image to project
// POST /api/projects/:id/images
func (h *ProjectHandler) AddImage(c *gin.Context) {
	var uri struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		ImageURL    string `json:"image_url" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.projectClient.AddProjectImage(ctx, &pb.AddProjectImageRequest{
		ProjectId:   uri.ID,
		ImageUrl:    req.ImageURL,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Image)
}

// AddLink adds link to project
// POST /api/projects/:id/links
func (h *ProjectHandler) AddLink(c *gin.Context) {
	var uri struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		LinkURL  string `json:"link_url" binding:"required"`
		LinkType string `json:"link_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.projectClient.AddProjectLink(ctx, &pb.AddProjectLinkRequest{
		ProjectId: uri.ID,
		LinkUrl:   req.LinkURL,
		LinkType:  req.LinkType,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Link)
}

// ListSkills returns all skills
// GET /api/skills
func (h *ProjectHandler) ListSkills(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.projectClient.ListSkills(ctx, &pb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Skills)
}

// CreateSkill creates a new skill
// POST /api/skills
func (h *ProjectHandler) CreateSkill(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.projectClient.CreateSkill(ctx, &pb.CreateSkillRequest{Name: req.Name})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Skill)
}

// AddMember adds a member to project (MOCK)
// POST /api/projects/:id/members
func (h *ProjectHandler) AddMember(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		UserID int64  `json:"userId" binding:"required"`
		Role   string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock response
	c.JSON(http.StatusOK, gin.H{
		"message":    "Add member endpoint (Mock)",
		"project_id": id,
		"user_id":    req.UserID,
		"role":       req.Role,
	})
}

// RemoveMember removes a member from project (MOCK)
// DELETE /api/projects/:id/members/:memberId
func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectID := c.Param("id")
	memberID := c.Param("memberId")

	// Mock response
	c.JSON(http.StatusOK, gin.H{
		"message":    "Remove member endpoint (Mock)",
		"project_id": projectID,
		"member_id":  memberID,
	})
}
