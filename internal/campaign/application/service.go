package application

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/campaign/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/timeuuid"
)

type Service struct{ repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

func (s *Service) ListCampaigns() ([]domain.Campaign, error) { return s.repo.ListCampaigns() }
func (s *Service) ListLeads() ([]domain.Lead, error)         { return s.repo.ListLeads() }
func (s *Service) ListRoutes(campaignID string) ([]domain.CampaignRoute, error) {
	return s.repo.ListRoutes(campaignID)
}

type CreateCampaignRequest struct {
	Name        string `json:"name"`
	BranchID    string `json:"branchId"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	TargetLeads int    `json:"targetLeads"`
}

type CreateRouteRequest struct {
	Date       string             `json:"date"`
	AssignedTo *string            `json:"assignedTo"`
	Status     string             `json:"status"`
	Stops      []domain.RouteStop `json:"stops"`
}

type CreateLeadRequest struct {
	FullName        string `json:"fullName"`
	Phone           string `json:"phone"`
	Suburb          string `json:"suburb"`
	CapturedBy      string `json:"capturedBy"`
	EstimatedAmount *int64 `json:"estimatedAmount"`
}

func (s *Service) CreateCampaign(req CreateCampaignRequest) (*domain.Campaign, error) {
	item := domain.Campaign{
		ID:          timeuuid.NewString(),
		Name:        req.Name,
		BranchID:    req.BranchID,
		Status:      "planning",
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		TargetLeads: req.TargetLeads,
	}
	return &item, s.repo.CreateCampaign(item)
}

func (s *Service) CreateRoute(campaignID string, req CreateRouteRequest) (*domain.CampaignRoute, error) {
	item := domain.CampaignRoute{
		ID:         timeuuid.NewString(),
		CampaignID: campaignID,
		Date:       req.Date,
		AssignedTo: req.AssignedTo,
		Status:     req.Status,
		Stops:      req.Stops,
	}
	return &item, s.repo.CreateRoute(item)
}

func (s *Service) CreateLead(campaignID string, req CreateLeadRequest) (*domain.Lead, error) {
	item := domain.Lead{
		ID:              timeuuid.NewString(),
		CampaignID:      campaignID,
		FullName:        req.FullName,
		Phone:           req.Phone,
		Suburb:          req.Suburb,
		CapturedBy:      req.CapturedBy,
		CapturedAt:      time.Now().UTC().Format(time.RFC3339),
		Qualified:       "pending",
		EstimatedAmount: req.EstimatedAmount,
	}
	return &item, s.repo.CreateLead(item)
}

func (s *Service) QualifyLead(id, status string) (*domain.Lead, error) {
	item, err := s.repo.GetLead(id)
	if err != nil || item == nil {
		return nil, httpx.NewError(fiber.StatusNotFound, "LEAD_NOT_FOUND", "lead not found")
	}
	item.Qualified = status
	return item, s.repo.UpdateLead(*item)
}

func EncodeStops(stops []domain.RouteStop) string {
	raw, _ := json.Marshal(stops)
	return string(raw)
}
