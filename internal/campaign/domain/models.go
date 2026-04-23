package domain

type Campaign struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	BranchID      string `json:"branchId"`
	Status        string `json:"status"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	TargetLeads   int    `json:"targetLeads"`
	CapturedLeads int    `json:"capturedLeads"`
}

type RouteStop struct {
	ID        string  `json:"id"`
	Address   string  `json:"address"`
	Suburb    string  `json:"suburb"`
	Status    string  `json:"status"`
	VisitedAt *string `json:"visitedAt,omitempty"`
	LeadID    *string `json:"leadId,omitempty"`
}

type CampaignRoute struct {
	ID         string      `json:"id"`
	CampaignID string      `json:"campaignId"`
	Date       string      `json:"date"`
	AssignedTo *string     `json:"assignedTo,omitempty"`
	Status     string      `json:"status"`
	Stops      []RouteStop `json:"stops"`
}

type Lead struct {
	ID              string `json:"id"`
	CampaignID      string `json:"campaignId"`
	FullName        string `json:"fullName"`
	Phone           string `json:"phone"`
	Suburb          string `json:"suburb"`
	CapturedBy      string `json:"capturedBy"`
	CapturedAt      string `json:"capturedAt"`
	Qualified       string `json:"qualified"`
	EstimatedAmount *int64 `json:"estimatedAmount,omitempty"`
}

type Repository interface {
	ListCampaigns() ([]Campaign, error)
	GetCampaign(id string) (*Campaign, error)
	CreateCampaign(campaign Campaign) error
	ListRoutes(campaignID string) ([]CampaignRoute, error)
	CreateRoute(route CampaignRoute) error
	ListLeads() ([]Lead, error)
	CreateLead(lead Lead) error
	UpdateLead(lead Lead) error
	GetLead(id string) (*Lead, error)
}
