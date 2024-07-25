package model

type Document struct {
	Id                  int          `json:"id"`
	Type                string       `json:"type"`
	FriendlyDescription string       `json:"friendlyDescription"`
	CreatedDate         string       `json:"createdDate"`
	Direction           string       `json:"direction"`
	Filename            string       `json:"filename"`
	CreatedBy           User         `json:"createdBy"`
	ReceivedDateTime    string       `json:"receivedDateTime"`
	Note                DocumentNote `json:"note"`
	ReformattedTime     string
}

type DocumentNote struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}
