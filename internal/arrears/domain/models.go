package domain

type PTP struct {
	ID         string  `json:"id"`
	Amount     int64   `json:"amount"`
	PromisedAt string  `json:"promisedAt"`
	CapturedBy string  `json:"capturedBy"`
	Note       *string `json:"note,omitempty"`
	Status     string  `json:"status"`
}

type ArrearsCase struct {
	ID                 string  `json:"id"`
	LoanID             string  `json:"loanId"`
	ClientName         string  `json:"clientName"`
	ClientPhone        string  `json:"clientPhone"`
	BranchID           string  `json:"branchId"`
	DaysPastDue        int     `json:"daysPastDue"`
	ArrearsAmount      int64   `json:"arrearsAmount"`
	OutstandingBalance int64   `json:"outstandingBalance"`
	Status             string  `json:"status"`
	AssignedTo         *string `json:"assignedTo,omitempty"`
	LastActionAt       *string `json:"lastActionAt,omitempty"`
	PTPs               []PTP   `json:"ptps"`
}

type ArrearsPayment struct {
	ID           string  `json:"id"`
	CaseID       string  `json:"caseId"`
	Amount       int64   `json:"amount"`
	Method       string  `json:"method"`
	Reference    string  `json:"reference"`
	CapturedBy   string  `json:"capturedBy"`
	CapturedAt   string  `json:"capturedAt"`
	AttachmentID *string `json:"attachmentId,omitempty"`
}

type Repository interface {
	ListCases() ([]ArrearsCase, error)
	GetCase(id string) (*ArrearsCase, error)
	UpsertCase(item ArrearsCase) error
	CreatePayment(payment ArrearsPayment) error
}
