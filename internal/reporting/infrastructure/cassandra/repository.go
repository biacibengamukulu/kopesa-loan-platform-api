package cassandra

import "github.com/gocql/gocql"

type Repository struct{ session *gocql.Session }

func NewRepository(session *gocql.Session) *Repository { return &Repository{session: session} }

func (r *Repository) GetExecOverview(period string) (map[string]any, error) {
	return map[string]any{
		"period": period,
		"loans": map[string]any{
			"disbursedAmount": 450000000,
			"disbursedCount":  312,
			"activeCount":     1840,
		},
		"arrears": map[string]any{
			"totalArrears": 87500000,
			"casesOpen":    412,
			"ptpRate":      0.41,
		},
		"campaigns": map[string]any{
			"active":         3,
			"leadsCaptured":  125,
			"leadsQualified": 71,
			"conversionRate": 0.18,
		},
		"trend": []map[string]any{
			{"date": "2026-04-01", "disbursed": 15000000, "collected": 12000000},
		},
	}, nil
}
