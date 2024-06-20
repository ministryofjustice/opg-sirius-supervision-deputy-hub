package model

type Document struct {
	Id                  int    `json:"id"`
	Type                string `json:"type"`
	FriendlyDescription string `json:"friendlyDescription"`
	CreatedDate         string `json:"createdDate"`
	Direction           string `json:"direction"`
	Filename            string `json:"filename"`
	CreatedBy           struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Email       string `json:"email"`
		Surname     string `json:"surname"`
	} `json:"createdBy"`
	ReceivedDateTime string `json:"receivedDateTime"`
}
