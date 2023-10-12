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
