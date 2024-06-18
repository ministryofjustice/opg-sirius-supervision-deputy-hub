package model

type Document struct {
	DisplayDate         string        `json:"displayDate"`
	Id                  int           `json:"id"`
	Uuid                string        `json:"uuid"`
	Type                string        `json:"type"`
	FriendlyDescription string        `json:"friendlyDescription"`
	CreatedDate         string        `json:"createdDate"`
	Direction           string        `json:"direction"`
	Filename            string        `json:"filename"`
	MimeType            string        `json:"mimeType"`
	CaseItems           []interface{} `json:"caseItems"`
	Persons             []struct {
		UId string `json:"uId"`
	} `json:"persons"`
	CreatedBy struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Email       string `json:"email"`
		Surname     string `json:"surname"`
	} `json:"createdBy"`
	ReceivedDateTime string `json:"receivedDateTime"`
	DocumentSource   string `json:"documentSource"`
	Metadata         struct {
		Type                string `json:"type"`
		Year                int    `json:"year"`
		ReportId            string `json:"report_id"`
		SubmissionId        int    `json:"submission_id"`
		ReportingPeriodTo   string `json:"reporting_period_to"`
		ReportingPeriodFrom string `json:"reporting_period_from"`
	} `json:"metadata"`
	ChildCount  int    `json:"childCount"`
	Subtype     string `json:"subtype"`
	Description string `json:"description"`
}
