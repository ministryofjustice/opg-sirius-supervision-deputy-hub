package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateAssuranceVisit struct {
	AssuranceType string `json:"assuranceType"`
	RequestedDate string `json:"requestedDate"`
	RequestedBy   int    `json:"requestedBy"`
}

func (c *Client) AddAssuranceVisit(ctx Context, assuranceType string, requestedDate string, userId, deputyId int) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(CreateAssuranceVisit{
		AssuranceType: assuranceType,
		RequestedDate: requestedDate,
		RequestedBy:   userId,
	})

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/deputies/%d/assurance-visits", deputyId), &body)

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{Errors: v.ValidationErrors}
		}

		return newStatusError(resp)
	}

	return nil
}
