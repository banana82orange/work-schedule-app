package entity

import "time"

// MediaFile represents a media file entity
type MediaFile struct {
	ID         int64     `json:"id"`
	FileName   string    `json:"file_name"`
	FileURL    string    `json:"file_url"`
	UploadedBy int64     `json:"uploaded_by"`
	UploadedAt time.Time `json:"uploaded_at"`
	FileType   string    `json:"file_type"` // image, document, resume
	FileSize   int64     `json:"file_size"`
}

// NewMediaFile creates a new media file entity
func NewMediaFile(fileName, fileURL, fileType string, uploadedBy, fileSize int64) *MediaFile {
	return &MediaFile{
		FileName:   fileName,
		FileURL:    fileURL,
		UploadedBy: uploadedBy,
		UploadedAt: time.Now(),
		FileType:   fileType,
		FileSize:   fileSize,
	}
}

// File type constants
const (
	FileTypeImage    = "image"
	FileTypeDocument = "document"
	FileTypeResume   = "resume"
)

// ValidFileTypes returns all valid file types
func ValidFileTypes() []string {
	return []string{FileTypeImage, FileTypeDocument, FileTypeResume}
}

// IsValidFileType checks if file type is valid
func IsValidFileType(fileType string) bool {
	for _, t := range ValidFileTypes() {
		if t == fileType {
			return true
		}
	}
	return false
}
