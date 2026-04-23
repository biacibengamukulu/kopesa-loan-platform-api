package application

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/loan/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/timeuuid"
)

type Service struct{ repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

type CreateApplicationRequest struct {
	ClientName      string `json:"clientName"`
	ClientID        string `json:"clientId"`
	ClientPhone     string `json:"clientPhone"`
	MonthlyIncome   int64  `json:"monthlyIncome"`
	MonthlyExpenses int64  `json:"monthlyExpenses"`
	RequestedAmount int64  `json:"requestedAmount"`
	TermMonths      int    `json:"termMonths"`
	BranchID        string `json:"branchId"`
	CreatedBy       string `json:"createdBy"`
}

type AssessApplicationRequest struct {
	Affordability int64  `json:"affordability"`
	Note          string `json:"note"`
}

type ApproveApplicationRequest struct {
	ApprovalNote string `json:"approvalNote"`
	ApprovedBy   string `json:"approvedBy"`
}

type DisburseLoanRequest struct {
	DisbursementRef string `json:"disbursementRef"`
	DisbursedAt     string `json:"disbursedAt"`
	DisbursedBy     string `json:"disbursedBy"`
}

func (s *Service) ListLoans(status, branchID string) ([]domain.Loan, error) {
	return s.repo.ListLoans(status, branchID)
}

func (s *Service) GetLoan(id string) (*domain.Loan, error) {
	loan, err := s.repo.GetLoan(id)
	if err != nil || loan == nil {
		return nil, httpx.NewError(fiber.StatusNotFound, "LOAN_NOT_FOUND", "loan not found")
	}
	return loan, nil
}

func (s *Service) ListApplications(status string) ([]domain.LoanApplication, error) {
	return s.repo.ListApplications(status)
}

func (s *Service) GetApplication(id string) (*domain.LoanApplication, error) {
	app, err := s.repo.GetApplication(id)
	if err != nil || app == nil {
		return nil, httpx.NewError(fiber.StatusNotFound, "LOAN_APPLICATION_NOT_FOUND", "application not found")
	}
	return app, nil
}

func (s *Service) CreateApplication(req CreateApplicationRequest) (*domain.LoanApplication, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	app := domain.LoanApplication{
		ID:              timeuuid.NewString(),
		LoanID:          fmt.Sprintf("LN-%d", time.Now().Unix()%100000),
		ClientName:      req.ClientName,
		ClientID:        req.ClientID,
		ClientPhone:     req.ClientPhone,
		MonthlyIncome:   req.MonthlyIncome,
		MonthlyExpenses: req.MonthlyExpenses,
		RequestedAmount: req.RequestedAmount,
		TermMonths:      req.TermMonths,
		BranchID:        req.BranchID,
		CreatedBy:       req.CreatedBy,
		Status:          "draft",
	}
	_ = now
	if err := s.repo.CreateApplication(app); err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *Service) AssessApplication(id string, req AssessApplicationRequest) (*domain.LoanApplication, error) {
	app, err := s.GetApplication(id)
	if err != nil {
		return nil, err
	}
	app.Status = "assessed"
	app.Affordability = &req.Affordability
	app.AssessmentNote = &req.Note
	if err := s.repo.UpdateApplication(*app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *Service) ApproveApplication(id string, req ApproveApplicationRequest) (*domain.Loan, error) {
	app, err := s.GetApplication(id)
	if err != nil {
		return nil, err
	}
	app.Status = "approved"
	app.ApprovalNote = &req.ApprovalNote
	if err := s.repo.UpdateApplication(*app); err != nil {
		return nil, err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	loan := domain.Loan{
		ID:                 timeuuid.NewString(),
		Reference:          app.LoanID,
		ClientName:         app.ClientName,
		ClientPhone:        app.ClientPhone,
		BranchID:           app.BranchID,
		Principal:          app.RequestedAmount,
		TermMonths:         app.TermMonths,
		RateBP:             2800,
		Status:             "approved",
		OutstandingBalance: app.RequestedAmount,
		CreatedBy:          app.CreatedBy,
		ApprovedBy:         &req.ApprovedBy,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := s.repo.CreateLoan(loan); err != nil {
		return nil, err
	}
	return &loan, nil
}

func (s *Service) DisburseLoan(id string, req DisburseLoanRequest) (*domain.Loan, error) {
	loan, err := s.GetLoan(id)
	if err != nil {
		return nil, err
	}
	loan.Status = "disbursed"
	loan.DisbursedBy = &req.DisbursedBy
	loan.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.repo.UpdateLoan(*loan); err != nil {
		return nil, err
	}
	return loan, nil
}
