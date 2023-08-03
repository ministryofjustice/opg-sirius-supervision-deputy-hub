package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type completeTask struct {
	Notes         string `json:"taskCompletedNotes"`
	CompletedById int    `json:"completedBy"`
}

func (c *Client) CompleteTask(ctx Context, userId, taskId int, notes string) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(completeTask{
		Notes:         notes,
		CompletedById: userId,
	})

	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/tasks/%d/mark-as-completed", taskId)

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
