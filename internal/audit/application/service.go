package application

import (
	"time"

	"github.com/biangacila/kopesa-loan-platform-api/internal/audit/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/timeuuid"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type RecordInput struct {
	Actor      string
	ActorRole  string
	Action     string
	EntityType string
	EntityID   string
	IP         string
	UserAgent  string
	RequestID  string
	Metadata   map[string]string
}

func (s *Service) Record(input RecordInput) error {
	return s.repo.Create(domain.Event{
		ID:         timeuuid.NewString(),
		OccurredAt: time.Now().UTC().Format(time.RFC3339),
		Actor:      input.Actor,
		ActorRole:  input.ActorRole,
		Action:     input.Action,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		IP:         input.IP,
		UserAgent:  input.UserAgent,
		RequestID:  input.RequestID,
		Metadata:   input.Metadata,
	})
}

func (s *Service) List(limit int) ([]domain.Event, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.repo.List(limit)
}
