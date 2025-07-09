package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AssignAssuranceVisitToClientsParams struct {
	DueDate   string   `json:"dueDate"`
	ClientIds []string `json:"clientIds"`
}

type ReassignResponse struct {
	ReassignName string `json:"reassignName"`
}

func (c *Client) AssignAssuranceVisitToClients(ctx Context, params AssignAssuranceVisitToClientsParams, deputyId int) (string, error) {
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

	if resp.StatusCode == http.StatusUnauthorized {
		return "", ErrUnauthorized
	}

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
