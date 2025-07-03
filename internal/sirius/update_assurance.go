package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateAssuranceDetails struct {
	CommissionedDate   string `json:"commissionedDate"`
	ReportDueDate      string `json:"reportDueDate"`
	ReportReceivedDate string `json:"reportReceivedDate"`
	VisitOutcome       string `json:"assuranceVisitOutcome"`
	PdrOutcome         string `json:"pdrOutcome"`
	ReportReviewDate   string `json:"reportReviewDate"`
	ReportMarkedAs     string `json:"reportMarkedAs"`
	VisitorAllocated   string `json:"visitorAllocated"`
	ReviewedBy         int    `json:"reviewedBy"`
	Note               string `json:"note"`
}

func (c *Client) UpdateAssurance(ctx Context, form UpdateAssuranceDetails, deputyId, visitId int) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(form)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/assurances/%d", deputyId, visitId), &body)

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
