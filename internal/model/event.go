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
	Assignee              string    `json:"assignee"`
	Changes               []Changes `json:"changes"`
	Client                []Client  `json:"additionalPersons"`
	Contact               Contact   `json:"deputyContact"`
	ClientCount           string    `json:"clientCount"`
	DueDate               string    `json:"dueDate"`
	DeputyID              string    `json:"personId"`
	DeputyName            string    `json:"personName"`
	Description           string    `json:"description"`
	Direction             string    `json:"direction"`
	ExecutiveCaseManager  string    `json:"executiveCaseManager"`
	Filename              string    `json:"filename"`
	Notes                 string    `json:"notes"`
	OrderType             string    `json:"orderType"`
	OrganisationName      string    `json:"organisationName"`
	OrderNumber           string    `json:"orderCourtRef"`
	OldAssigneeName       string    `json:"oldAssigneeName"`
	ReceivedDate          string    `json:"receivedDate"`
	RecipientEmailAddress string    `json:"recipientEmailAddress"`
	SiriusId              string    `json:"orderUid"`
	TaskType              string    `json:"taskType"`
	TaskCompletedNotes    string    `json:"taskCompletedNotes"`
	Type                  string    `json:"type"`
}

type Changes struct {
	FieldName string `json:"fieldName"`
	OldValue  string `json:"oldValue"`
	NewValue  string `json:"newValue"`
}
