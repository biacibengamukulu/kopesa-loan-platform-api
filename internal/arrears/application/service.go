package application

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/arrears/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/timeuuid"
)

type Service struct{ repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

type AllocateRequest struct {
	AssignedTo string `json:"assignedTo"`
}

type PTPRequest struct {
	Amount     int64  `json:"amount"`
	PromisedAt string `json:"promisedAt"`
	Note       string `json:"note"`
	CapturedBy string `json:"capturedBy"`
}

type PaymentRequest struct {
	Amount       int64   `json:"amount"`
	Method       string  `json:"method"`
	Reference    string  `json:"reference"`
	CapturedBy   string  `json:"capturedBy"`
	CapturedAt   string  `json:"capturedAt"`
	AttachmentID *string `json:"attachmentId"`
}

func (s *Service) ListCases() ([]domain.ArrearsCase, error) { return s.repo.ListCases() }

func (s *Service) GetCase(id string) (*domain.ArrearsCase, error) {
	item, err := s.repo.GetCase(id)
	if err != nil || item == nil {
		return nil, httpx.NewError(fiber.StatusNotFound, "ARREARS_NOT_FOUND", "arrears case not found")
	}
	return item, nil
}

func (s *Service) Allocate(id string, req AllocateRequest) (*domain.ArrearsCase, error) {
	item, err := s.GetCase(id)
	if err != nil {
		return nil, err
	}
	item.Status = "allocated"
	item.AssignedTo = &req.AssignedTo
	now := time.Now().UTC().Format(time.RFC3339)
	item.LastActionAt = &now
	return item, s.repo.UpsertCase(*item)
}

func (s *Service) CreatePTP(id string, req PTPRequest) (*domain.ArrearsCase, error) {
	item, err := s.GetCase(id)
	if err != nil {
		return nil, err
	}
	note := req.Note
	item.PTPs = append(item.PTPs, domain.PTP{
		ID:         timeuuid.NewString(),
		Amount:     req.Amount,
		PromisedAt: req.PromisedAt,
		CapturedBy: req.CapturedBy,
		Note:       &note,
		Status:     "pending",
	})
	item.Status = "ptp"
	now := time.Now().UTC().Format(time.RFC3339)
	item.LastActionAt = &now
	return item, s.repo.UpsertCase(*item)
}

func (s *Service) CapturePayment(id string, req PaymentRequest) (map[string]any, error) {
	item, err := s.GetCase(id)
	if err != nil {
		return nil, err
	}
	payment := domain.ArrearsPayment{
		ID:           timeuuid.NewString(),
		CaseID:       id,
		Amount:       req.Amount,
		Method:       req.Method,
		Reference:    req.Reference,
		CapturedBy:   req.CapturedBy,
		CapturedAt:   req.CapturedAt,
		AttachmentID: req.AttachmentID,
	}
	item.ArrearsAmount -= req.Amount
	if item.ArrearsAmount <= 0 {
		item.ArrearsAmount = 0
		item.Status = "paid"
	}
	now := time.Now().UTC().Format(time.RFC3339)
	item.LastActionAt = &now
	if err := s.repo.CreatePayment(payment); err != nil {
		return nil, err
	}
	if err := s.repo.UpsertCase(*item); err != nil {
		return nil, err
	}
	return map[string]any{"case": item, "payment": payment}, nil
}
