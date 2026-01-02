package repository

import (
	"context"

	"github.com/portfolio/media-service/internal/domain/entity"
)

// MediaFileRepository defines the interface for media file data access
type MediaFileRepository interface {
	Create(ctx context.Context, file *entity.MediaFile) error
	GetByID(ctx context.Context, id int64) (*entity.MediaFile, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int, fileType string) ([]*entity.MediaFile, int, error)
	GetByUserID(ctx context.Context, userID int64, page, limit int) ([]*entity.MediaFile, int, error)
}

// FileStorage defines the interface for file storage operations
type FileStorage interface {
	Save(ctx context.Context, fileName string, data []byte) (string, error)
	Delete(ctx context.Context, fileURL string) error
	Get(ctx context.Context, fileURL string) ([]byte, error)
}
