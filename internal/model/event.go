package model

type DeputyEvent struct {
	ID         int    `json:"id"`
	Timestamp  string `json:"timestamp"`
	EventType  string `json:"eventType"`
	User       User   `json:"user"`
	Event      Event  `json:"event"`
	IsNewEvent bool
}

type Event struct {
	OrderType            string    `json:"orderType"`
	SiriusId             string    `json:"orderUid"`
	OrderNumber          string    `json:"orderCourtRef"`
	DeputyID             string    `json:"personId"`
	DeputyName           string    `json:"personName"`
	OrganisationName     string    `json:"organisationName"`
	ExecutiveCaseManager string    `json:"executiveCaseManager"`
	Changes              []Changes `json:"changes"`
	Client               []Client  `json:"additionalPersons"`
	RequestedBy          string    `json:"requestedBy"`
	RequestedDate        string    `json:"requestedDate"`
	CommissionedDate     string    `json:"commissionedDate"`
	ReportDueDate        string    `json:"reportDueDate"`
	ReportReceivedDate   string    `json:"reportReceivedDate"`
	VisitOutcome         string    `json:"assuranceVisitOutcome"`
	ReportReviewDate     string    `json:"reportReviewDate"`
	VisitReportMarkedAs  string    `json:"assuranceVisitReportMarkedAs"`
	VisitorAllocated     string    `json:"visitorAllocated"`
	ReviewedBy           string    `json:"reviewedBy"`
	Contact              Contact   `json:"deputyContact"`
	TaskType             string    `json:"taskType"`
	Assignee             string    `json:"assignee"`
	DueDate              string    `json:"dueDate"`
	Notes                string    `json:"description"`
	OldAssigneeName      string    `json:"oldAssigneeName"`
	TaskCompletedNotes   string    `json:"taskCompletedNotes"`
}

type Changes struct {
	FieldName string `json:"fieldName"`
	OldValue  string `json:"oldValue"`
	NewValue  string `json:"newValue"`
}
