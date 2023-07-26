package model

type TeamMember struct {
	ID          int
	DisplayName string
	CurrentEcm  int
}

type Team struct {
	Members []TeamMember
	Name    string
}
