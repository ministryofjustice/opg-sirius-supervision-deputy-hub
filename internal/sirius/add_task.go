package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type addTask struct {
	DeputyId    int    `json:"personId"`
	TaskType    string `json:"type"`
	DueDate     string `json:"dueDate"`
	Notes       string `json:"notes"`
	IsCaseOwner bool   `json:"isCaseOwner"` // temporary until assignee selection is added
}

func (c *Client) AddTask(ctx Context, deputyId int, taskType string, dueDate string, notes string) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(addTask{
		TaskType:    taskType,
		DueDate:     FormatDateTime(IsoDateTimeZone, dueDate, SiriusDate),
		Notes:       notes,
		IsCaseOwner: true,
		DeputyId:    deputyId,
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
