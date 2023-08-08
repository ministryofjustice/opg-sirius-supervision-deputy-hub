package model

type Client struct {
	ID       string `json:"personId"`
	Uid      string `json:"personUid"`
	Name     string `json:"personName"`
	CourtRef string `json:"personCourtRef"`
}
