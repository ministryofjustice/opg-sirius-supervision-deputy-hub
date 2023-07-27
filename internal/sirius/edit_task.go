package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type editTask struct {
	CaseOwnerTask bool   `json:"isCaseOwner"`
	DueDate       string `json:"dueDate"`
	Notes         string `json:"description"`
	AssigneeId    int    `json:"assigneeId"`
	DeputyId      int    `json:"personId"`
	IsDeputyTask  bool   `json:"isDeputyTask"`
}

func (c *Client) EditTask(ctx Context, deputyId int, taskId int, dueDate string, notes string, assigneeId int) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(editTask{
		DueDate:       FormatDateTime(IsoDate, dueDate, SiriusDate),
		Notes:         notes,
		AssigneeId:    assigneeId,
		CaseOwnerTask: false,
		DeputyId:      deputyId,
		IsDeputyTask:  true,
	})

	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/tasks/%d", taskId)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

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

	if resp.StatusCode != http.StatusOK {
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
