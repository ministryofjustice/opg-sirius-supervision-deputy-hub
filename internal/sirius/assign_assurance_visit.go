package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AssignAssuranceVisitToClientsParams struct {
	DueDate   string
	ClientIds []string `json:"clientIds"`
}

type ReassignResponse struct {
	ReassignName string `json:"reassignName"`
}

func (c *Client) AssignAssuranceVisitToClients(ctx Context, params AssignAssuranceVisitToClientsParams, deputyId int) (string, error) {
	var u ReassignResponse
	var body bytes.Buffer
	var err error

	err = json.NewEncoder(&body).Encode(params)

	if err != nil {
		return "", err
	}
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d", deputyId), &body)

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

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return "", &ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return "", newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&u)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("You have assigned %d clients for an assurance visit", len(params.ClientIds)), nil
}
