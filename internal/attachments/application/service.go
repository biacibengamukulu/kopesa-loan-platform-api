package application

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path"
	"time"

	"github.com/biangacila/kopesa-loan-platform-api/internal/attachments/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/timeuuid"
)

type Service struct {
	repo    domain.Repository
	storage domain.Storage
}

func NewService(repo domain.Repository, storage domain.Storage) *Service {
	return &Service{repo: repo, storage: storage}
}

type PresignRequest struct {
	FileName  string `json:"fileName"`
	MimeType  string `json:"mimeType"`
	SizeBytes int64  `json:"sizeBytes"`
	Context   string `json:"context"`
	EntityID  string `json:"entityId"`
}

type FinalizeRequest struct {
	ID         string  `json:"id"`
	Context    string  `json:"context"`
	EntityID   string  `json:"entityId"`
	FileName   string  `json:"fileName"`
	MimeType   string  `json:"mimeType"`
	SizeBytes  int64   `json:"sizeBytes"`
	URL        string  `json:"url"`
	CapturedBy string  `json:"capturedBy"`
	Note       *string `json:"note"`
}

func (s *Service) List(entityID, context string) ([]domain.Attachment, error) {
	return s.repo.List(entityID, context)
}

func (s *Service) Get(id string) (*domain.Attachment, error) {
	item, err := s.repo.Get(id)
	if err != nil || item == nil {
		return nil, httpx.NewError(http.StatusNotFound, "ATTACHMENT_NOT_FOUND", "attachment not found")
	}
	return item, nil
}

func (s *Service) Presign(req PresignRequest) map[string]any {
	id := timeuuid.NewString()
	return map[string]any{
		"attachmentId": id,
		"uploadUrl":    fmt.Sprintf("/v1/attachments/upload?attachmentId=%s&context=%s&entityId=%s&fileName=%s", id, req.Context, req.EntityID, req.FileName),
		"expiresIn":    600,
	}
}

func (s *Service) Finalize(req FinalizeRequest) (*domain.Attachment, error) {
	id := req.ID
	if id == "" {
		id = timeuuid.NewString()
	}
	item := domain.Attachment{
		ID:         id,
		Context:    req.Context,
		EntityID:   req.EntityID,
		FileName:   req.FileName,
		MimeType:   req.MimeType,
		SizeBytes:  req.SizeBytes,
		URL:        req.URL,
		CapturedBy: req.CapturedBy,
		CapturedAt: time.Now().UTC().Format(time.RFC3339),
		Sync:       "pending",
		Note:       req.Note,
		Provider:   "dropbox",
	}
	return &item, s.repo.Create(item)
}

type UploadDirectRequest struct {
	ID         string
	Context    string
	EntityID   string
	FileName   string
	MimeType   string
	SizeBytes  int64
	CapturedBy string
	Note       *string
	File       multipart.File
}

func (s *Service) Upload(req UploadDirectRequest) (*domain.Attachment, error) {
	id := req.ID
	if id == "" {
		id = timeuuid.NewString()
	}
	targetPath := path.Join(req.Context, req.EntityID, req.FileName)
	uploaded, err := s.storage.Upload(context.Background(), domain.UploadRequest{
		Path:     targetPath,
		FileName: req.FileName,
		File:     req.File,
	})
	if err != nil {
		return nil, err
	}
	item := domain.Attachment{
		ID:         id,
		Context:    req.Context,
		EntityID:   req.EntityID,
		FileName:   req.FileName,
		MimeType:   req.MimeType,
		SizeBytes:  req.SizeBytes,
		URL:        s.storage.DownloadURL(uploaded.Revision),
		CapturedBy: req.CapturedBy,
		CapturedAt: time.Now().UTC().Format(time.RFC3339),
		Sync:       "synced",
		Note:       req.Note,
		Provider:   "dropbox",
		Path:       uploaded.Path,
		Revision:   uploaded.Revision,
	}
	return &item, s.repo.Create(item)
}
