package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BulkAssignAssuranceVisitTasksToClientsParams struct {
	DueDate   string   `json:"dueDate"`
	ClientIds []string `json:"clientIds"`
}

func (c *Client) BulkAssignAssuranceVisitTasksToClients(ctx Context, params BulkAssignAssuranceVisitTasksToClientsParams, deputyId int) (string, error) {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(params)
	if err != nil {
		return "", err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/bulk-assurance-visit-tasks", deputyId), &body)

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode >= 300 {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return "", ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return "", newStatusError(resp)
	}

	return fmt.Sprintf("You have assigned %d clients for an assurance visit", len(params.ClientIds)), nil
}
