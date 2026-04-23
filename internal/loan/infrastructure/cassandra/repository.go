package cassandra

import (
	"github.com/gocql/gocql"

	"github.com/biangacila/kopesa-loan-platform-api/internal/loan/domain"
)

type Repository struct{ session *gocql.Session }

func NewRepository(session *gocql.Session) *Repository { return &Repository{session: session} }

func (r *Repository) ListLoans(status, branchID string) ([]domain.Loan, error) {
	query := `SELECT id, reference, client_name, client_phone, branch_id, principal, term_months, rate_bp, status, outstanding_balance, next_due_date, created_by, assessed_by, approved_by, disbursed_by, created_at, updated_at FROM loans_loans`
	if status != "" || branchID != "" {
		query += " ALLOW FILTERING"
	}
	iter := r.session.Query(query).Iter()
	defer iter.Close()
	out := make([]domain.Loan, 0)
	for {
		var item domain.Loan
		if !iter.Scan(&item.ID, &item.Reference, &item.ClientName, &item.ClientPhone, &item.BranchID, &item.Principal, &item.TermMonths, &item.RateBP, &item.Status, &item.OutstandingBalance, &item.NextDueDate, &item.CreatedBy, &item.AssessedBy, &item.ApprovedBy, &item.DisbursedBy, &item.CreatedAt, &item.UpdatedAt) {
			break
		}
		if status != "" && item.Status != status {
			continue
		}
		if branchID != "" && item.BranchID != branchID {
			continue
		}
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) GetLoan(id string) (*domain.Loan, error) {
	iter := r.session.Query(`SELECT id, reference, client_name, client_phone, branch_id, principal, term_months, rate_bp, status, outstanding_balance, next_due_date, created_by, assessed_by, approved_by, disbursed_by, created_at, updated_at FROM loans_loans WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var item domain.Loan
	if !iter.Scan(&item.ID, &item.Reference, &item.ClientName, &item.ClientPhone, &item.BranchID, &item.Principal, &item.TermMonths, &item.RateBP, &item.Status, &item.OutstandingBalance, &item.NextDueDate, &item.CreatedBy, &item.AssessedBy, &item.ApprovedBy, &item.DisbursedBy, &item.CreatedAt, &item.UpdatedAt) {
		return nil, nil
	}
	return &item, iter.Close()
}

func (r *Repository) ListApplications(status string) ([]domain.LoanApplication, error) {
	iter := r.session.Query(`SELECT id, loan_id, client_name, client_id, client_phone, monthly_income, monthly_expenses, requested_amount, term_months, branch_id, created_by, status, affordability, assessment_note, approval_note FROM loans_applications`).Iter()
	defer iter.Close()
	out := make([]domain.LoanApplication, 0)
	for {
		var item domain.LoanApplication
		if !iter.Scan(&item.ID, &item.LoanID, &item.ClientName, &item.ClientID, &item.ClientPhone, &item.MonthlyIncome, &item.MonthlyExpenses, &item.RequestedAmount, &item.TermMonths, &item.BranchID, &item.CreatedBy, &item.Status, &item.Affordability, &item.AssessmentNote, &item.ApprovalNote) {
			break
		}
		if status != "" && item.Status != status {
			continue
		}
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) GetApplication(id string) (*domain.LoanApplication, error) {
	iter := r.session.Query(`SELECT id, loan_id, client_name, client_id, client_phone, monthly_income, monthly_expenses, requested_amount, term_months, branch_id, created_by, status, affordability, assessment_note, approval_note FROM loans_applications WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var item domain.LoanApplication
	if !iter.Scan(&item.ID, &item.LoanID, &item.ClientName, &item.ClientID, &item.ClientPhone, &item.MonthlyIncome, &item.MonthlyExpenses, &item.RequestedAmount, &item.TermMonths, &item.BranchID, &item.CreatedBy, &item.Status, &item.Affordability, &item.AssessmentNote, &item.ApprovalNote) {
		return nil, nil
	}
	return &item, iter.Close()
}

func (r *Repository) CreateApplication(app domain.LoanApplication) error {
	return r.session.Query(`INSERT INTO loans_applications (id, loan_id, client_name, client_id, client_phone, monthly_income, monthly_expenses, requested_amount, term_months, branch_id, created_by, status, affordability, assessment_note, approval_note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		app.ID, app.LoanID, app.ClientName, app.ClientID, app.ClientPhone, app.MonthlyIncome, app.MonthlyExpenses, app.RequestedAmount, app.TermMonths, app.BranchID, app.CreatedBy, app.Status, app.Affordability, app.AssessmentNote, app.ApprovalNote,
	).Exec()
}

func (r *Repository) UpdateApplication(app domain.LoanApplication) error {
	return r.CreateApplication(app)
}

func (r *Repository) CreateLoan(loan domain.Loan) error {
	return r.session.Query(`INSERT INTO loans_loans (id, reference, client_name, client_phone, branch_id, principal, term_months, rate_bp, status, outstanding_balance, next_due_date, created_by, assessed_by, approved_by, disbursed_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		loan.ID, loan.Reference, loan.ClientName, loan.ClientPhone, loan.BranchID, loan.Principal, loan.TermMonths, loan.RateBP, loan.Status, loan.OutstandingBalance, loan.NextDueDate, loan.CreatedBy, loan.AssessedBy, loan.ApprovedBy, loan.DisbursedBy, loan.CreatedAt, loan.UpdatedAt,
	).Exec()
}

func (r *Repository) UpdateLoan(loan domain.Loan) error { return r.CreateLoan(loan) }
