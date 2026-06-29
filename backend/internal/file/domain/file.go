package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// File holds metadata for an uploaded file. The bytes live on disk under StoredName.
type File struct {
	ID           uuid.UUID `json:"id"`
	OriginalName string    `json:"name"`
	StoredName   string    `json:"-"`
	Path         string    `json:"-"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime"`
	UploadedBy   uuid.UUID `json:"uploadedBy"`
	CreatedAt    time.Time `json:"createdAt"`
}

// AllowedExtensions is the accepted upload allowlist (lowercase, no dot).
var AllowedExtensions = map[string]bool{
	// documents
	"pdf": true, "doc": true, "docx": true, "xls": true, "xlsx": true,
	"ppt": true, "pptx": true, "txt": true, "csv": true,
	// archives
	"zip": true, "rar": true, "7z": true,
	// images
	"png": true, "jpg": true, "jpeg": true, "gif": true, "webp": true, "svg": true,
	// video
	"mp4": true, "avi": true, "mov": true, "mkv": true,
	// audio
	"mp3": true, "wav": true,
}

// ExtAllowed reports whether an extension (with or without leading dot) is accepted.
func ExtAllowed(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	return AllowedExtensions[ext]
}
