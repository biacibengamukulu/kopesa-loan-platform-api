package application

import "github.com/biangacila/kopesa-loan-platform-api/internal/reporting/domain"

type Service struct{ repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

func (s *Service) ExecOverview(period string) (map[string]any, error) {
	if period == "" {
		period = "mtd"
	}
	return s.repo.GetExecOverview(period)
}
