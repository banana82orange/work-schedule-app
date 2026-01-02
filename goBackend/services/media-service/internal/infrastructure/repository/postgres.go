package repository

import (
	"context"
	"database/sql"

	"github.com/portfolio/media-service/internal/domain/entity"
)

// PostgresMediaFileRepository implements MediaFileRepository
type PostgresMediaFileRepository struct {
	db *sql.DB
}

// NewPostgresMediaFileRepository creates a new repository
func NewPostgresMediaFileRepository(db *sql.DB) *PostgresMediaFileRepository {
	return &PostgresMediaFileRepository{db: db}
}

// Create creates a new media file record
func (r *PostgresMediaFileRepository) Create(ctx context.Context, file *entity.MediaFile) error {
	query := `
		INSERT INTO media_files (file_name, file_url, uploaded_by, uploaded_at, file_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		file.FileName, file.FileURL, file.UploadedBy, file.UploadedAt, file.FileType,
	).Scan(&file.ID)
}

// GetByID gets a media file by ID
func (r *PostgresMediaFileRepository) GetByID(ctx context.Context, id int64) (*entity.MediaFile, error) {
	query := `SELECT id, file_name, file_url, uploaded_by, uploaded_at, file_type FROM media_files WHERE id = $1`
	file := &entity.MediaFile{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&file.ID, &file.FileName, &file.FileURL, &file.UploadedBy, &file.UploadedAt, &file.FileType,
	)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Delete deletes a media file record
func (r *PostgresMediaFileRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM media_files WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List lists media files with pagination
func (r *PostgresMediaFileRepository) List(ctx context.Context, page, limit int, fileType string) ([]*entity.MediaFile, int, error) {
	offset := (page - 1) * limit

	// Build query
	var countQuery, query string
	var args []interface{}

	if fileType != "" {
		countQuery = `SELECT COUNT(*) FROM media_files WHERE file_type = $1`
		query = `SELECT id, file_name, file_url, uploaded_by, uploaded_at, file_type FROM media_files WHERE file_type = $1 ORDER BY uploaded_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{fileType, limit, offset}
	} else {
		countQuery = `SELECT COUNT(*) FROM media_files`
		query = `SELECT id, file_name, file_url, uploaded_by, uploaded_at, file_type FROM media_files ORDER BY uploaded_at DESC LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	}

	// Get total
	var total int
	if fileType != "" {
		if err := r.db.QueryRowContext(ctx, countQuery, fileType).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	// Get files
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var files []*entity.MediaFile
	for rows.Next() {
		file := &entity.MediaFile{}
		if err := rows.Scan(&file.ID, &file.FileName, &file.FileURL, &file.UploadedBy, &file.UploadedAt, &file.FileType); err != nil {
			return nil, 0, err
		}
		files = append(files, file)
	}

	return files, total, nil
}

// GetByUserID gets files uploaded by a user
func (r *PostgresMediaFileRepository) GetByUserID(ctx context.Context, userID int64, page, limit int) ([]*entity.MediaFile, int, error) {
	offset := (page - 1) * limit

	// Get total
	var total int
	countQuery := `SELECT COUNT(*) FROM media_files WHERE uploaded_by = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get files
	query := `SELECT id, file_name, file_url, uploaded_by, uploaded_at, file_type FROM media_files WHERE uploaded_by = $1 ORDER BY uploaded_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var files []*entity.MediaFile
	for rows.Next() {
		file := &entity.MediaFile{}
		if err := rows.Scan(&file.ID, &file.FileName, &file.FileURL, &file.UploadedBy, &file.UploadedAt, &file.FileType); err != nil {
			return nil, 0, err
		}
		files = append(files, file)
	}

	return files, total, nil
}
