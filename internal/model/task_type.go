package model

type TaskType struct {
	Handle        string `json:"handle"`
	Description   string `json:"incomplete"`
	ProDeputyTask bool   `json:"proDeputyTask"`
	PaDeputyTask  bool   `json:"paDeputyTask"`
}
