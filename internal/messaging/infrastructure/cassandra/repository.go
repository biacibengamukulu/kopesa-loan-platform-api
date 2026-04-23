package cassandra

import (
	"github.com/gocql/gocql"

	"github.com/biangacila/kopesa-loan-platform-api/internal/messaging/domain"
)

type Repository struct{ session *gocql.Session }

func NewRepository(session *gocql.Session) *Repository { return &Repository{session: session} }

func (r *Repository) ListTemplates() ([]domain.MessageTemplate, error) {
	iter := r.session.Query(`SELECT id, name, channel, context, body, description FROM messaging_templates`).Iter()
	defer iter.Close()
	out := make([]domain.MessageTemplate, 0)
	for {
		var item domain.MessageTemplate
		if !iter.Scan(&item.ID, &item.Name, &item.Channel, &item.Context, &item.Body, &item.Description) {
			break
		}
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) GetTemplate(id string) (*domain.MessageTemplate, error) {
	iter := r.session.Query(`SELECT id, name, channel, context, body, description FROM messaging_templates WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var item domain.MessageTemplate
	if !iter.Scan(&item.ID, &item.Name, &item.Channel, &item.Context, &item.Body, &item.Description) {
		return nil, nil
	}
	return &item, iter.Close()
}

func (r *Repository) ListLogs(context, entityID string) ([]domain.MessageLogEntry, error) {
	iter := r.session.Query(`SELECT id, context, entity_id, channel, template_id, recipient_to, body, status, sent_at, sent_by, next_touch_at, provider_ref FROM messaging_logs`).Iter()
	defer iter.Close()
	out := make([]domain.MessageLogEntry, 0)
	for {
		var item domain.MessageLogEntry
		if !iter.Scan(&item.ID, &item.Context, &item.EntityID, &item.Channel, &item.TemplateID, &item.To, &item.Body, &item.Status, &item.SentAt, &item.SentBy, &item.NextTouchAt, &item.ProviderRef) {
			break
		}
		if context != "" && item.Context != context {
			continue
		}
		if entityID != "" && item.EntityID != entityID {
			continue
		}
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) CreateLog(entry domain.MessageLogEntry) error {
	return r.session.Query(`INSERT INTO messaging_logs (id, context, entity_id, channel, template_id, recipient_to, body, status, sent_at, sent_by, next_touch_at, provider_ref) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.Context, entry.EntityID, entry.Channel, entry.TemplateID, entry.To, entry.Body, entry.Status, entry.SentAt, entry.SentBy, entry.NextTouchAt, entry.ProviderRef,
	).Exec()
}
