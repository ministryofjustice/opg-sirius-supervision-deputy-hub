package model

type DeputyEvent struct {
	ID         int    `json:"id"`
	Timestamp  string `json:"timestamp"`
	EventType  string `json:"eventType"`
	User       User   `json:"user"`
	Event      Event  `json:"event"`
	IsNewEvent bool
}

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"displayName"`
	PhoneNumber string `json:"phoneNumber"`
}

type Event struct {
	OrderType            string         `json:"orderType"`
	SiriusId             string         `json:"orderUid"`
	OrderNumber          string         `json:"orderCourtRef"`
	DeputyID             string         `json:"personId"`
	DeputyName           string         `json:"personName"`
	OrganisationName     string         `json:"organisationName"`
	ExecutiveCaseManager string         `json:"executiveCaseManager"`
	Changes              []Changes      `json:"changes"`
	Client               []ClientPerson `json:"additionalPersons"`
	RequestedBy          string         `json:"requestedBy"`
	RequestedDate        string         `json:"requestedDate"`
	CommissionedDate     string         `json:"commissionedDate"`
	ReportDueDate        string         `json:"reportDueDate"`
	ReportReceivedDate   string         `json:"reportReceivedDate"`
	VisitOutcome         string         `json:"assuranceVisitOutcome"`
	ReportReviewDate     string         `json:"reportReviewDate"`
	VisitReportMarkedAs  string         `json:"assuranceVisitReportMarkedAs"`
	VisitorAllocated     string         `json:"visitorAllocated"`
	ReviewedBy           string         `json:"reviewedBy"`
	Contact              ContactPerson  `json:"deputyContact"`
	TaskType             string         `json:"taskType"`
	Assignee             string         `json:"assignee"`
	DueDate              string         `json:"dueDate"`
	Notes                string         `json:"description"`
	OldAssigneeName      string         `json:"oldAssigneeName"`
	TaskCompletedNotes   string         `json:"taskCompletedNotes"`
}

type Changes struct {
	FieldName string `json:"fieldName"`
	OldValue  string `json:"oldValue"`
	NewValue  string `json:"newValue"`
}

type ClientPerson struct {
	ID       string `json:"personId"`
	Uid      string `json:"personUid"`
	Name     string `json:"personName"`
	CourtRef string `json:"personCourtRef"`
}

type ContactPerson struct {
	Name             string `json:"name"`
	JobTitle         string `json:"jobTitle"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	OtherPhoneNumber string `json:"otherPhoneNumber"`
	Notes            string `json:"notes"`
}
