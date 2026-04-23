package domain

import (
	"context"
	"mime/multipart"
)

type Attachment struct {
	ID         string  `json:"id"`
	Context    string  `json:"context"`
	EntityID   string  `json:"entityId"`
	FileName   string  `json:"fileName"`
	MimeType   string  `json:"mimeType"`
	SizeBytes  int64   `json:"sizeBytes"`
	URL        string  `json:"url"`
	CapturedBy string  `json:"capturedBy"`
	CapturedAt string  `json:"capturedAt"`
	Sync       string  `json:"sync"`
	Note       *string `json:"note,omitempty"`
	Provider   string  `json:"provider,omitempty"`
	Path       string  `json:"path,omitempty"`
	Revision   string  `json:"revision,omitempty"`
}

type Repository interface {
	List(entityID, context string) ([]Attachment, error)
	Get(id string) (*Attachment, error)
	Create(item Attachment) error
}

type UploadRequest struct {
	Path     string
	FileName string
	File     multipart.File
}

type UploadResult struct {
	Path     string
	Revision string
}

type Storage interface {
	Upload(ctx context.Context, req UploadRequest) (*UploadResult, error)
	DownloadURL(revision string) string
}
