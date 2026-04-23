package domain

type ExecOverview struct {
	Period    string           `json:"period"`
	Loans     map[string]any   `json:"loans"`
	Arrears   map[string]any   `json:"arrears"`
	Campaigns map[string]any   `json:"campaigns"`
	Trend     []map[string]any `json:"trend"`
}

type Repository interface {
	GetExecOverview(period string) (map[string]any, error)
}
