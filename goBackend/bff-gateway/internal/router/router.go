package router

import (
	"github.com/gin-gonic/gin"
	"github.com/portfolio/bff-gateway/internal/grpc"
	"github.com/portfolio/bff-gateway/internal/handler"
	"github.com/portfolio/bff-gateway/internal/middleware"
)

// SetupRouter configures all routes
func SetupRouter(jwtSecret string, clients *grpc.ClientManager) *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(clients.GetAuthConn())
	projectHandler := handler.NewProjectHandler(clients.GetProjectConn())
	taskHandler := handler.NewTaskHandler(clients.GetTaskConn())
	analyticsHandler := handler.NewAnalyticsHandler(clients.GetAnalyticsConn())
	mediaHandler := handler.NewMediaHandler(clients.GetMediaConn())

	// ==========================================
	// Auth routes (public)
	// ==========================================
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/validate", authHandler.ValidateToken)
	}

	// ==========================================
	// Protected routes (require authentication)
	// ==========================================
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Auth - Profile
		protected.GET("/auth/profile", authHandler.GetProfile)

		// Users (admin only)
		users := protected.Group("/users")
		users.Use(middleware.RoleMiddleware("admin"))
		{
			users.GET("", authHandler.ListUsers)
			users.GET("/:id", authHandler.GetUser)
			users.PUT("/:id", authHandler.UpdateUser)
			users.DELETE("/:id", authHandler.DeleteUser)
		}

		// ==========================================
		// Projects
		// ==========================================
		projects := protected.Group("/projects")
		{
			projects.POST("", projectHandler.CreateProject)
			projects.GET("", projectHandler.ListProjects)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)

			// Project skills
			projects.POST("/:id/skills", projectHandler.AddSkill)

			// Project tech
			projects.POST("/:id/tech", projectHandler.AddTech)

			// Project images
			projects.POST("/:id/images", projectHandler.AddImage)

			// Project links
			projects.POST("/:id/links", projectHandler.AddLink)

			// Project members
			projects.POST("/:id/members", projectHandler.AddMember)
			projects.DELETE("/:id/members/:memberId", projectHandler.RemoveMember)
		}

		// Skills
		skills := protected.Group("/skills")
		{
			skills.GET("", projectHandler.ListSkills)
			skills.POST("", projectHandler.CreateSkill)
		}

		// ==========================================
		// Tasks
		// ==========================================
		tasks := protected.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.ListTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)

			// Subtasks
			tasks.POST("/:id/subtasks", taskHandler.CreateSubtask)
			tasks.GET("/:id/subtasks", taskHandler.ListSubtasks)

			// Comments
			tasks.POST("/:id/comments", taskHandler.AddComment)
			tasks.GET("/:id/comments", taskHandler.ListComments)

			// Attachments
			tasks.POST("/:id/attachments", taskHandler.AddAttachment)
			tasks.GET("/:id/attachments", taskHandler.ListAttachments)

			// Tags
			tasks.POST("/:id/tags", taskHandler.AddTag)
		}

		// Tags
		tags := protected.Group("/tags")
		{
			tags.GET("", taskHandler.ListTags)
			tags.POST("", taskHandler.CreateTag)
		}

		// ==========================================
		// Analytics
		// ==========================================
		analytics := protected.Group("/analytics")
		{
			// Dashboard
			analytics.GET("/dashboard", analyticsHandler.GetDashboardStats)

			// Project analytics
			analytics.POST("/projects/:id/view", analyticsHandler.RecordProjectView)
			analytics.GET("/projects/:id/views", analyticsHandler.GetProjectViews)
			analytics.GET("/projects/:id/stats", analyticsHandler.GetProjectStats)

			// Task analytics
			analytics.POST("/tasks/:id/activity", analyticsHandler.RecordTaskActivity)
			analytics.GET("/tasks/:id/activities", analyticsHandler.GetTaskActivities)
		}

		// ==========================================
		// Media
		// ==========================================
		media := protected.Group("/media")
		{
			media.POST("/upload", mediaHandler.UploadFile)
			media.GET("", mediaHandler.ListFiles)
			media.GET("/my-files", mediaHandler.GetUserFiles)
			media.GET("/:id", mediaHandler.GetFile)
			media.DELETE("/:id", mediaHandler.DeleteFile)
		}
	}

	return r
}
