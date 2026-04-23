package domain

type Loan struct {
	ID                 string  `json:"id"`
	Reference          string  `json:"reference"`
	ClientName         string  `json:"clientName"`
	ClientPhone        string  `json:"clientPhone"`
	BranchID           string  `json:"branchId"`
	Principal          int64   `json:"principal"`
	TermMonths         int     `json:"termMonths"`
	RateBP             int     `json:"rateBp"`
	Status             string  `json:"status"`
	OutstandingBalance int64   `json:"outstandingBalance"`
	NextDueDate        *string `json:"nextDueDate,omitempty"`
	CreatedBy          string  `json:"createdBy"`
	AssessedBy         *string `json:"assessedBy,omitempty"`
	ApprovedBy         *string `json:"approvedBy,omitempty"`
	DisbursedBy        *string `json:"disbursedBy,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type LoanApplication struct {
	ID              string  `json:"id"`
	LoanID          string  `json:"loanId"`
	ClientName      string  `json:"clientName"`
	ClientID        string  `json:"clientId"`
	ClientPhone     string  `json:"clientPhone"`
	MonthlyIncome   int64   `json:"monthlyIncome"`
	MonthlyExpenses int64   `json:"monthlyExpenses"`
	RequestedAmount int64   `json:"requestedAmount"`
	TermMonths      int     `json:"termMonths"`
	BranchID        string  `json:"branchId"`
	CreatedBy       string  `json:"createdBy"`
	Status          string  `json:"status"`
	Affordability   *int64  `json:"affordability,omitempty"`
	AssessmentNote  *string `json:"assessmentNote,omitempty"`
	ApprovalNote    *string `json:"approvalNote,omitempty"`
}

type Repository interface {
	ListLoans(status, branchID string) ([]Loan, error)
	GetLoan(id string) (*Loan, error)
	ListApplications(status string) ([]LoanApplication, error)
	GetApplication(id string) (*LoanApplication, error)
	CreateApplication(app LoanApplication) error
	UpdateApplication(app LoanApplication) error
	CreateLoan(loan Loan) error
	UpdateLoan(loan Loan) error
}
