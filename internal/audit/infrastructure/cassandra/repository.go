package cassandra

import (
	"github.com/gocql/gocql"

	"github.com/biangacila/kopesa-loan-platform-api/internal/audit/domain"
)

type Repository struct {
	session *gocql.Session
}

func NewRepository(session *gocql.Session) *Repository {
	return &Repository{session: session}
}

func (r *Repository) Create(event domain.Event) error {
	return r.session.Query(`INSERT INTO audit_events (id, occurred_at, actor, actor_role, action, entity_type, entity_id, ip, user_agent, request_id, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.OccurredAt, event.Actor, event.ActorRole, event.Action, event.EntityType, event.EntityID, event.IP, event.UserAgent, event.RequestID, event.Metadata,
	).Exec()
}

func (r *Repository) List(limit int) ([]domain.Event, error) {
	iter := r.session.Query(`SELECT id, occurred_at, actor, actor_role, action, entity_type, entity_id, ip, user_agent, request_id, metadata FROM audit_events LIMIT ?`, limit).Iter()
	defer iter.Close()

	events := make([]domain.Event, 0, limit)
	for {
		var event domain.Event
		if !iter.Scan(&event.ID, &event.OccurredAt, &event.Actor, &event.ActorRole, &event.Action, &event.EntityType, &event.EntityID, &event.IP, &event.UserAgent, &event.RequestID, &event.Metadata) {
			break
		}
		events = append(events, event)
	}
	return events, iter.Close()
}
