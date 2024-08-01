package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type addTask struct {
	DeputyId   int    `json:"personId"`
	TaskType   string `json:"type"`
	TypeName   string `json:"name"`
	DueDate    string `json:"dueDate"`
	Notes      string `json:"description"`
	AssigneeId int    `json:"assigneeId"`
}

func (c *ApiClient) AddTask(ctx Context, deputyId int, taskType string, typeName string, dueDate string, notes string, assigneeId int) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(addTask{
		TaskType:   taskType,
		TypeName:   typeName,
		DueDate:    FormatDateTime(IsoDate, dueDate, SiriusDate),
		Notes:      notes,
		AssigneeId: assigneeId,
		DeputyId:   deputyId,
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/tasks", &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusCreated {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{Errors: v.ValidationErrors}
		}

		return newStatusError(resp)
	}

	return nil
}
