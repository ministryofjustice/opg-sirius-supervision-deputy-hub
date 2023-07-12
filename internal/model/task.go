package model

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type Assignee struct {
	Id          int           `json:"id"`
	Teams       []sirius.Team `json:"teams,omitempty"`
	DisplayName string        `json:"displayName"`
}

type Task struct {
	Id            int      `json:"id"`
	Type          string   `json:"type"`
	DueDate       string   `json:"dueDate"`
	Name          string   `json:"name"`
	Assignee      Assignee `json:"assignee"`
	CreatedTime   string   `json:"createdTime"`
	CaseOwnerTask bool     `json:"caseOwnerTask"`
	Notes         string   `json:"description"`
}

func (t Task) GetName(taskTypes []TaskType) string {
	for _, taskType := range taskTypes {
		if t.Type == taskType.Handle {
			return taskType.Description
		}
	}
	return t.Name
}
