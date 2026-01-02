package usecase

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/portfolio/media-service/internal/domain/entity"
	"github.com/portfolio/media-service/internal/domain/repository"
)

var (
	ErrFileNotFound    = errors.New("file not found")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrUploadFailed    = errors.New("upload failed")
)

// MediaUseCase handles media business logic
type MediaUseCase struct {
	fileRepo repository.MediaFileRepository
	storage  repository.FileStorage
}

// NewMediaUseCase creates a new MediaUseCase
func NewMediaUseCase(fileRepo repository.MediaFileRepository, storage repository.FileStorage) *MediaUseCase {
	return &MediaUseCase{
		fileRepo: fileRepo,
		storage:  storage,
	}
}

// UploadFile uploads a file
func (uc *MediaUseCase) UploadFile(ctx context.Context, fileName, fileType string, uploadedBy int64, data []byte) (*entity.MediaFile, error) {
	if !entity.IsValidFileType(fileType) {
		return nil, ErrInvalidFileType
	}

	// Generate unique filename
	ext := filepath.Ext(fileName)
	uniqueName := time.Now().Format("20060102150405") + "_" + fileName

	// Save to storage
	fileURL, err := uc.storage.Save(ctx, uniqueName, data)
	if err != nil {
		return nil, ErrUploadFailed
	}

	// Create file record
	file := entity.NewMediaFile(fileName, fileURL, fileType, uploadedBy, int64(len(data)))
	if ext != "" {
		file.FileName = fileName
	}

	if err := uc.fileRepo.Create(ctx, file); err != nil {
		// Cleanup uploaded file on error
		_ = uc.storage.Delete(ctx, fileURL)
		return nil, err
	}

	return file, nil
}

// GetFile retrieves a file by ID
func (uc *MediaUseCase) GetFile(ctx context.Context, id int64) (*entity.MediaFile, error) {
	file, err := uc.fileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

// DeleteFile deletes a file
func (uc *MediaUseCase) DeleteFile(ctx context.Context, id int64) error {
	file, err := uc.fileRepo.GetByID(ctx, id)
	if err != nil {
		return ErrFileNotFound
	}

	// Delete from storage
	if err := uc.storage.Delete(ctx, file.FileURL); err != nil {
		return err
	}

	// Delete record
	return uc.fileRepo.Delete(ctx, id)
}

// ListFiles lists files with pagination
func (uc *MediaUseCase) ListFiles(ctx context.Context, page, limit int, fileType string) ([]*entity.MediaFile, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.fileRepo.List(ctx, page, limit, fileType)
}

// GetFilesByUser gets files by user
func (uc *MediaUseCase) GetFilesByUser(ctx context.Context, userID int64, page, limit int) ([]*entity.MediaFile, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.fileRepo.GetByUserID(ctx, userID, page, limit)
}
