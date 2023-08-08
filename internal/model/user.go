package model

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"displayName"`
	PhoneNumber string `json:"phoneNumber"`
}
