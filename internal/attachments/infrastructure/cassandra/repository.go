package cassandra

import (
	"github.com/gocql/gocql"

	"github.com/biangacila/kopesa-loan-platform-api/internal/attachments/domain"
)

type Repository struct{ session *gocql.Session }

func NewRepository(session *gocql.Session) *Repository { return &Repository{session: session} }

func (r *Repository) List(entityID, context string) ([]domain.Attachment, error) {
	iter := r.session.Query(`SELECT id, context, entity_id, file_name, mime_type, size_bytes, url, captured_by, captured_at, sync, note, provider, path, revision FROM attachments_attachments`).Iter()
	defer iter.Close()
	out := make([]domain.Attachment, 0)
	for {
		var item domain.Attachment
		if !iter.Scan(&item.ID, &item.Context, &item.EntityID, &item.FileName, &item.MimeType, &item.SizeBytes, &item.URL, &item.CapturedBy, &item.CapturedAt, &item.Sync, &item.Note, &item.Provider, &item.Path, &item.Revision) {
			break
		}
		if entityID != "" && item.EntityID != entityID {
			continue
		}
		if context != "" && item.Context != context {
			continue
		}
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) Get(id string) (*domain.Attachment, error) {
	iter := r.session.Query(`SELECT id, context, entity_id, file_name, mime_type, size_bytes, url, captured_by, captured_at, sync, note, provider, path, revision FROM attachments_attachments WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var item domain.Attachment
	if !iter.Scan(&item.ID, &item.Context, &item.EntityID, &item.FileName, &item.MimeType, &item.SizeBytes, &item.URL, &item.CapturedBy, &item.CapturedAt, &item.Sync, &item.Note, &item.Provider, &item.Path, &item.Revision) {
		return nil, nil
	}
	return &item, iter.Close()
}

func (r *Repository) Create(item domain.Attachment) error {
	var note any
	if item.Note != nil {
		note = *item.Note
	}
	return r.session.Query(`INSERT INTO attachments_attachments (id, context, entity_id, file_name, mime_type, size_bytes, url, captured_by, captured_at, sync, note, provider, path, revision) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.Context, item.EntityID, item.FileName, item.MimeType, item.SizeBytes, item.URL, item.CapturedBy, item.CapturedAt, item.Sync, note, item.Provider, item.Path, item.Revision,
	).Exec()
}
