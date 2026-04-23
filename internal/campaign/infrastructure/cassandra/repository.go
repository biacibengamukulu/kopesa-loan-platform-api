package cassandra

import (
	"encoding/json"

	"github.com/gocql/gocql"

	"github.com/biangacila/kopesa-loan-platform-api/internal/campaign/domain"
)

type Repository struct{ session *gocql.Session }

func NewRepository(session *gocql.Session) *Repository { return &Repository{session: session} }

func (r *Repository) ListCampaigns() ([]domain.Campaign, error) {
	iter := r.session.Query(`SELECT id, name, branch_id, status, start_date, end_date, target_leads, captured_leads FROM campaigns_campaigns`).Iter()
	defer iter.Close()
	out := make([]domain.Campaign, 0)
	for {
		var item domain.Campaign
		if !iter.Scan(&item.ID, &item.Name, &item.BranchID, &item.Status, &item.StartDate, &item.EndDate, &item.TargetLeads, &item.CapturedLeads) {
			break
		}
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) GetCampaign(id string) (*domain.Campaign, error) {
	iter := r.session.Query(`SELECT id, name, branch_id, status, start_date, end_date, target_leads, captured_leads FROM campaigns_campaigns WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var item domain.Campaign
	if !iter.Scan(&item.ID, &item.Name, &item.BranchID, &item.Status, &item.StartDate, &item.EndDate, &item.TargetLeads, &item.CapturedLeads) {
		return nil, nil
	}
	return &item, iter.Close()
}

func (r *Repository) CreateCampaign(campaign domain.Campaign) error {
	return r.session.Query(`INSERT INTO campaigns_campaigns (id, name, branch_id, status, start_date, end_date, target_leads, captured_leads) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		campaign.ID, campaign.Name, campaign.BranchID, campaign.Status, campaign.StartDate, campaign.EndDate, campaign.TargetLeads, campaign.CapturedLeads,
	).Exec()
}

func (r *Repository) ListRoutes(campaignID string) ([]domain.CampaignRoute, error) {
	iter := r.session.Query(`SELECT id, campaign_id, date, assigned_to, status, stops_json FROM campaigns_routes WHERE campaign_id = ?`, campaignID).Iter()
	defer iter.Close()
	out := make([]domain.CampaignRoute, 0)
	for {
		var item domain.CampaignRoute
		var stopsJSON string
		if !iter.Scan(&item.ID, &item.CampaignID, &item.Date, &item.AssignedTo, &item.Status, &stopsJSON) {
			break
		}
		_ = json.Unmarshal([]byte(stopsJSON), &item.Stops)
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) CreateRoute(route domain.CampaignRoute) error {
	stopsJSON, _ := json.Marshal(route.Stops)
	return r.session.Query(`INSERT INTO campaigns_routes (campaign_id, id, date, assigned_to, status, stops_json) VALUES (?, ?, ?, ?, ?, ?)`,
		route.CampaignID, route.ID, route.Date, route.AssignedTo, route.Status, string(stopsJSON),
	).Exec()
}

func (r *Repository) ListLeads() ([]domain.Lead, error) {
	iter := r.session.Query(`SELECT id, campaign_id, full_name, phone, suburb, captured_by, captured_at, qualified, estimated_amount FROM campaigns_leads`).Iter()
	defer iter.Close()
	out := make([]domain.Lead, 0)
	for {
		var item domain.Lead
		if !iter.Scan(&item.ID, &item.CampaignID, &item.FullName, &item.Phone, &item.Suburb, &item.CapturedBy, &item.CapturedAt, &item.Qualified, &item.EstimatedAmount) {
			break
		}
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) CreateLead(lead domain.Lead) error {
	return r.session.Query(`INSERT INTO campaigns_leads (id, campaign_id, full_name, phone, suburb, captured_by, captured_at, qualified, estimated_amount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		lead.ID, lead.CampaignID, lead.FullName, lead.Phone, lead.Suburb, lead.CapturedBy, lead.CapturedAt, lead.Qualified, lead.EstimatedAmount,
	).Exec()
}

func (r *Repository) UpdateLead(lead domain.Lead) error { return r.CreateLead(lead) }

func (r *Repository) GetLead(id string) (*domain.Lead, error) {
	iter := r.session.Query(`SELECT id, campaign_id, full_name, phone, suburb, captured_by, captured_at, qualified, estimated_amount FROM campaigns_leads WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var item domain.Lead
	if !iter.Scan(&item.ID, &item.CampaignID, &item.FullName, &item.Phone, &item.Suburb, &item.CapturedBy, &item.CapturedAt, &item.Qualified, &item.EstimatedAmount) {
		return nil, nil
	}
	return &item, iter.Close()
}
