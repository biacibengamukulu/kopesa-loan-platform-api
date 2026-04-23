package application

import (
	"context"
	"strings"
	"time"

	"github.com/biangacila/kopesa-loan-platform-api/internal/messaging/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/apperr"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/timeuuid"
)

type Service struct {
	repo    domain.Repository
	gateway domain.Gateway
}

func NewService(repo domain.Repository, gateway domain.Gateway) *Service {
	return &Service{repo: repo, gateway: gateway}
}

type SendRequest struct {
	Context     string            `json:"context"`
	EntityID    string            `json:"entityId"`
	Channel     string            `json:"channel"`
	To          string            `json:"to"`
	Subject     string            `json:"subject"`
	TemplateID  *string           `json:"templateId"`
	Body        string            `json:"body"`
	Variables   map[string]string `json:"variables"`
	NextTouchAt *string           `json:"nextTouchAt"`
	SentBy      string            `json:"sentBy"`
}

func (s *Service) ListTemplates() ([]domain.MessageTemplate, error) { return s.repo.ListTemplates() }
func (s *Service) ListLogs(context, entityID string) ([]domain.MessageLogEntry, error) {
	return s.repo.ListLogs(context, entityID)
}

func (s *Service) Send(req SendRequest) (*domain.MessageLogEntry, error) {
	if req.Body == "" && req.TemplateID != nil {
		template, err := s.repo.GetTemplate(*req.TemplateID)
		if err != nil {
			return nil, err
		}
		if template == nil {
			return nil, apperr.New(404, "MESSAGE_TEMPLATE_NOT_FOUND", "message template not found")
		}
		req.Body = renderTemplate(template.Body, req.Variables)
		if req.Channel == "" {
			req.Channel = template.Channel
		}
	}
	if req.Channel == "" {
		req.Channel = "sms"
	}
	result, err := s.gateway.Send(context.Background(), domain.DispatchRequest{
		Channel: req.Channel,
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
	})
	if err != nil {
		return nil, err
	}
	providerRef := result.ProviderRef
	entry := domain.MessageLogEntry{
		ID:          timeuuid.NewString(),
		Context:     req.Context,
		EntityID:    req.EntityID,
		Channel:     req.Channel,
		TemplateID:  req.TemplateID,
		To:          req.To,
		Body:        req.Body,
		Status:      result.Status,
		SentAt:      time.Now().UTC().Format(time.RFC3339),
		SentBy:      req.SentBy,
		NextTouchAt: req.NextTouchAt,
		ProviderRef: &providerRef,
	}
	return &entry, s.repo.CreateLog(entry)
}

func renderTemplate(body string, variables map[string]string) string {
	rendered := body
	for key, value := range variables {
		rendered = strings.ReplaceAll(rendered, "{{"+key+"}}", value)
	}
	return rendered
}
