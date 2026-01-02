package handler

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/portfolio/proto/media"
	"google.golang.org/grpc"
)

// MediaHandler handles media endpoints
type MediaHandler struct {
	mediaClient pb.MediaServiceClient
}

// NewMediaHandler creates a new MediaHandler
func NewMediaHandler(conn *grpc.ClientConn) *MediaHandler {
	return &MediaHandler{
		mediaClient: pb.NewMediaServiceClient(conn),
	}
}

const (
	// MaxFileSize is 10MB
	MaxFileSize = 10 << 20
	// ChunkSize is 64KB for streaming
	ChunkSize = 64 * 1024
)

// UploadFile uploads a file
// POST /api/media/upload
func (h *MediaHandler) UploadFile(c *gin.Context) {
	// Limit body size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required: " + err.Error()})
		return
	}
	defer file.Close()

	fileType := c.PostForm("file_type")
	if fileType == "" {
		fileType = "document"
	}

	userIDVal, _ := c.Get("user_id")
	var userID int64
	if v, ok := userIDVal.(float64); ok {
		userID = int64(v)
	} else if v, ok := userIDVal.(int64); ok {
		userID = v
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute) // Longer timeout for upload
	defer cancel()

	stream, err := h.mediaClient.UploadFile(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start upload: " + err.Error()})
		return
	}

	// 1. Send Metadata
	req := &pb.UploadFileRequest{
		Data: &pb.UploadFileRequest_Metadata{
			Metadata: &pb.FileMetadata{
				FileName:   header.Filename,
				FileType:   fileType,
				UploadedBy: userID,
			},
		},
	}
	if err := stream.Send(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send metadata: " + err.Error()})
		return
	}

	// 2. Send Chunks
	buffer := make([]byte, ChunkSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file: " + err.Error()})
			return
		}

		req := &pb.UploadFileRequest{
			Data: &pb.UploadFileRequest_Chunk{
				Chunk: buffer[:n],
			},
		}
		if err := stream.Send(req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send chunk: " + err.Error()})
			return
		}
	}

	// 3. Close and Recv
	resp, err := stream.CloseAndRecv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.File)
}

// GetFile returns a file by ID
// GET /api/media/:id
func (h *MediaHandler) GetFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.mediaClient.GetFile(ctx, &pb.GetFileRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.File)
}

// DeleteFile deletes a file
// DELETE /api/media/:id
func (h *MediaHandler) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.mediaClient.DeleteFile(ctx, &pb.DeleteFileRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// ListFiles returns list of files
// GET /api/media
func (h *MediaHandler) ListFiles(c *gin.Context) {
	// page := c.DefaultQuery("page", "1")
	// limit := c.DefaultQuery("limit", "10")
	fileType := c.Query("file_type")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.mediaClient.ListFiles(ctx, &pb.ListFilesRequest{
		Page:     1,
		Limit:    100,
		FileType: fileType,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Files) // Note: Proto response wraps files in 'Files' field? Yes checked proto.
}

// GetUserFiles returns files uploaded by current user
// GET /api/media/my-files
func (h *MediaHandler) GetUserFiles(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	var userID int64
	if v, ok := userIDVal.(float64); ok {
		userID = int64(v)
	} else if v, ok := userIDVal.(int64); ok {
		userID = v
	}
	// page := c.DefaultQuery("page", "1")
	// limit := c.DefaultQuery("limit", "10")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.mediaClient.GetFilesByUser(ctx, &pb.GetFilesByUserRequest{
		UserId: userID,
		Page:   1,
		Limit:  100,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Files)
}
