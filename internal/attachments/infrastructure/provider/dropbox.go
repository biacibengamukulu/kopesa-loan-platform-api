package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/biangacila/kopesa-loan-platform-api/internal/attachments/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/apperr"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/config"
)

type Dropbox struct {
	client  *http.Client
	baseURL string
}

func NewDropbox(cfg config.Config) *Dropbox {
	return &Dropbox{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: strings.TrimRight(cfg.DropboxBaseURL, "/"),
	}
}

func (d *Dropbox) Upload(ctx context.Context, req domain.UploadRequest) (*domain.UploadResult, error) {
	if seeker, ok := req.File.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, apperr.New(http.StatusBadRequest, "ATTACHMENT_FILE_INVALID", err.Error())
		}
	}
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, req.File); err != nil {
		return nil, err
	}
	_ = writer.Close()

	u := d.baseURL + "/upload?path=" + url.QueryEscape(req.Path)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, u, &body)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := d.client.Do(httpReq)
	if err != nil {
		return nil, apperr.New(http.StatusBadGateway, "DROPBOX_PROVIDER_DOWN", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, apperr.New(http.StatusBadGateway, "DROPBOX_UPLOAD_FAILED", string(raw))
	}
	rev, _ := d.findRevision(ctx, req.Path)
	return &domain.UploadResult{Path: req.Path, Revision: rev}, nil
}

func (d *Dropbox) DownloadURL(revision string) string {
	if revision == "" {
		return ""
	}
	return d.baseURL + "/stream/" + revision
}

func (d *Dropbox) findRevision(ctx context.Context, fullPath string) (string, error) {
	dir := path.Dir(fullPath)
	name := path.Base(fullPath)
	u := d.baseURL + "/list?path=" + url.QueryEscape(dir)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	resp, err := d.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", nil
	}
	var listing struct {
		Entries []struct {
			Name string `json:"name"`
			Rev  string `json:"rev"`
		} `json:"entries"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
		return "", nil
	}
	for _, entry := range listing.Entries {
		if entry.Name == name {
			return entry.Rev, nil
		}
	}
	return "", nil
}
