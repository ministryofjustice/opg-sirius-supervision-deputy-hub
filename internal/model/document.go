package model

type Document struct {
	Id                  int    `json:"id"`
	Type                string `json:"type"`
	FriendlyDescription string `json:"friendlyDescription"`
	CreatedDate         string `json:"createdDate"`
	Direction           string `json:"direction"`
	Filename            string `json:"filename"`
	CreatedBy           User   `json:"createdBy"`
	ReceivedDateTime    string `json:"receivedDateTime"`
}

type Document2 struct {
	Id                  int    `json:"id"`
	Type                string `json:"type"`
	FriendlyDescription string `json:"friendlyDescription"`
	CreatedDate         string `json:"createdDate"`
	Direction           string `json:"direction"`
	Filename            string `json:"filename"`
	ReceivedDateTime    string `json:"receivedDateTime"`
	Note                struct {
		Description string `json:"description"`
		Name        string `json:"name"`
	} `json:"note"`
}
