package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateGcmIssue struct {
	ClientCaseRecNumber string `json:"caseRecNumber"`
	GcmIssueType        string `json:"gcmIssueType"`
	Notes               string `json:"notes"`
}

func (c *Client) AddGcmIssue(ctx Context, clientCaseRecNumber, notes string, gcmIssueType string, deputyId int) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(CreateGcmIssue{
		ClientCaseRecNumber: clientCaseRecNumber,
		GcmIssueType:        gcmIssueType,
		Notes:               notes,
	})

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/gcm-issues", deputyId), &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

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
