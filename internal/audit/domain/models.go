package domain

type Event struct {
	ID         string            `json:"id"`
	OccurredAt string            `json:"occurredAt"`
	Actor      string            `json:"actor"`
	ActorRole  string            `json:"actorRole"`
	Action     string            `json:"action"`
	EntityType string            `json:"entityType"`
	EntityID   string            `json:"entityId"`
	IP         string            `json:"ip"`
	UserAgent  string            `json:"userAgent"`
	RequestID  string            `json:"requestId"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Repository interface {
	Create(event Event) error
	List(limit int) ([]Event, error)
}
