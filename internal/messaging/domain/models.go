package domain

import "context"

type MessageTemplate struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Channel     string  `json:"channel"`
	Context     string  `json:"context"`
	Body        string  `json:"body"`
	Description *string `json:"description,omitempty"`
}

type MessageLogEntry struct {
	ID          string  `json:"id"`
	Context     string  `json:"context"`
	EntityID    string  `json:"entityId"`
	Channel     string  `json:"channel"`
	TemplateID  *string `json:"templateId,omitempty"`
	To          string  `json:"to"`
	Body        string  `json:"body"`
	Status      string  `json:"status"`
	SentAt      string  `json:"sentAt"`
	SentBy      string  `json:"sentBy"`
	NextTouchAt *string `json:"nextTouchAt,omitempty"`
	ProviderRef *string `json:"providerRef,omitempty"`
}

type Repository interface {
	ListTemplates() ([]MessageTemplate, error)
	GetTemplate(id string) (*MessageTemplate, error)
	ListLogs(context, entityID string) ([]MessageLogEntry, error)
	CreateLog(entry MessageLogEntry) error
}

type DispatchRequest struct {
	Channel string
	To      string
	Subject string
	Body    string
}

type DispatchResult struct {
	Status      string
	ProviderRef string
}

type Gateway interface {
	Send(ctx context.Context, req DispatchRequest) (*DispatchResult, error)
}
