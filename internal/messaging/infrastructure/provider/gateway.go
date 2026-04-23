package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/biangacila/kopesa-loan-platform-api/internal/messaging/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/apperr"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/config"
)

type Gateway struct {
	client *http.Client
	cfg    config.Config
}

func NewGateway(cfg config.Config) *Gateway {
	return &Gateway{
		client: &http.Client{Timeout: 20 * time.Second},
		cfg:    cfg,
	}
}

func (g *Gateway) Send(ctx context.Context, req domain.DispatchRequest) (*domain.DispatchResult, error) {
	switch req.Channel {
	case "sms":
		return g.sendSMS(ctx, req)
	case "whatsapp":
		return g.sendWhatsApp(ctx, req)
	case "email":
		return g.sendEmail(ctx, req)
	case "both":
		if _, err := g.sendSMS(ctx, req); err != nil {
			return nil, err
		}
		return g.sendWhatsApp(ctx, req)
	default:
		return nil, apperr.New(http.StatusBadRequest, "MESSAGING_CHANNEL_UNSUPPORTED", "unsupported messaging channel")
	}
}

func (g *Gateway) sendSMS(ctx context.Context, req domain.DispatchRequest) (*domain.DispatchResult, error) {
	payload, _ := json.Marshal(map[string]string{
		"phone":   normalizePhone(req.To),
		"message": req.Body,
	})
	url := strings.TrimRight(g.cfg.ProviderBaseURL, "/") + "/send-sms/post"
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, apperr.New(http.StatusBadGateway, "SMS_PROVIDER_DOWN", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, apperr.New(http.StatusBadGateway, "SMS_PROVIDER_ERROR", string(body))
	}
	return &domain.DispatchResult{Status: "sent", ProviderRef: "sms:" + normalizePhone(req.To)}, nil
}

func (g *Gateway) sendEmail(ctx context.Context, req domain.DispatchRequest) (*domain.DispatchResult, error) {
	receivers := splitReceivers(req.To)
	payload, _ := json.Marshal(map[string]any{
		"from":     g.cfg.EmailFrom,
		"receiver": receivers,
		"subject":  req.Subject,
		"html":     req.Body,
		"status":   "PENDING",
		"retries":  0,
	})
	url := strings.TrimRight(g.cfg.ProviderBaseURL, "/") + "/send-email"
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, apperr.New(http.StatusBadGateway, "EMAIL_PROVIDER_DOWN", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, apperr.New(http.StatusBadGateway, "EMAIL_PROVIDER_ERROR", string(body))
	}
	return &domain.DispatchResult{Status: "queued", ProviderRef: "email:" + strings.Join(receivers, ",")}, nil
}

func (g *Gateway) sendWhatsApp(ctx context.Context, req domain.DispatchRequest) (*domain.DispatchResult, error) {
	payload, _ := json.Marshal(map[string]string{
		"number": normalizePhone(req.To),
		"text":   req.Body,
	})
	url := strings.TrimRight(g.cfg.WhatsAppBaseURL, "/") + "/message/sendText/" + g.cfg.WhatsAppInstance
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", g.cfg.WhatsAppAPIKey)
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, apperr.New(http.StatusBadGateway, "WHATSAPP_PROVIDER_DOWN", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, apperr.New(http.StatusBadGateway, "WHATSAPP_PROVIDER_ERROR", string(body))
	}
	return &domain.DispatchResult{Status: "sent", ProviderRef: "whatsapp:" + normalizePhone(req.To)}, nil
}

func normalizePhone(value string) string {
	value = strings.TrimSpace(value)
	return strings.TrimPrefix(value, "+")
}

func splitReceivers(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ';' })
	out := make([]string, 0, len(parts))
	for _, item := range parts {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	if len(out) == 0 {
		return []string{strings.TrimSpace(value)}
	}
	return out
}
