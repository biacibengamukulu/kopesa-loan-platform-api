package cassandra

import (
	"encoding/json"

	"github.com/gocql/gocql"

	"github.com/biangacila/kopesa-loan-platform-api/internal/arrears/domain"
)

type Repository struct{ session *gocql.Session }

func NewRepository(session *gocql.Session) *Repository { return &Repository{session: session} }

func (r *Repository) ListCases() ([]domain.ArrearsCase, error) {
	iter := r.session.Query(`SELECT id, loan_id, client_name, client_phone, branch_id, days_past_due, arrears_amount, outstanding_balance, status, assigned_to, last_action_at, ptps_json FROM arrears_cases`).Iter()
	defer iter.Close()
	out := make([]domain.ArrearsCase, 0)
	for {
		var item domain.ArrearsCase
		var ptpsJSON string
		if !iter.Scan(&item.ID, &item.LoanID, &item.ClientName, &item.ClientPhone, &item.BranchID, &item.DaysPastDue, &item.ArrearsAmount, &item.OutstandingBalance, &item.Status, &item.AssignedTo, &item.LastActionAt, &ptpsJSON) {
			break
		}
		_ = json.Unmarshal([]byte(ptpsJSON), &item.PTPs)
		out = append(out, item)
	}
	return out, iter.Close()
}

func (r *Repository) GetCase(id string) (*domain.ArrearsCase, error) {
	iter := r.session.Query(`SELECT id, loan_id, client_name, client_phone, branch_id, days_past_due, arrears_amount, outstanding_balance, status, assigned_to, last_action_at, ptps_json FROM arrears_cases WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var item domain.ArrearsCase
	var ptpsJSON string
	if !iter.Scan(&item.ID, &item.LoanID, &item.ClientName, &item.ClientPhone, &item.BranchID, &item.DaysPastDue, &item.ArrearsAmount, &item.OutstandingBalance, &item.Status, &item.AssignedTo, &item.LastActionAt, &ptpsJSON) {
		return nil, nil
	}
	_ = json.Unmarshal([]byte(ptpsJSON), &item.PTPs)
	return &item, iter.Close()
}

func (r *Repository) UpsertCase(item domain.ArrearsCase) error {
	ptpsJSON, _ := json.Marshal(item.PTPs)
	return r.session.Query(`INSERT INTO arrears_cases (id, loan_id, client_name, client_phone, branch_id, days_past_due, arrears_amount, outstanding_balance, status, assigned_to, last_action_at, ptps_json) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.LoanID, item.ClientName, item.ClientPhone, item.BranchID, item.DaysPastDue, item.ArrearsAmount, item.OutstandingBalance, item.Status, item.AssignedTo, item.LastActionAt, string(ptpsJSON),
	).Exec()
}

func (r *Repository) CreatePayment(payment domain.ArrearsPayment) error {
	return r.session.Query(`INSERT INTO arrears_payments (id, case_id, amount, method, reference, captured_by, captured_at, attachment_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		payment.ID, payment.CaseID, payment.Amount, payment.Method, payment.Reference, payment.CapturedBy, payment.CapturedAt, payment.AttachmentID,
	).Exec()
}
