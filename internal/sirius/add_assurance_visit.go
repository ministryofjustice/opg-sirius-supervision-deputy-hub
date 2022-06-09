package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AssuranceVisit struct {
	ID           int    `json:"id"`
	RequestedDate     string `json:"requstedDate"`
	RequestedById int `json:"requestedById"`
	CommissionedDate string `json:"commissionedDate"`
	Visitor string `json:"visitor"`
	ReportDueDate         string `json:"reportDueDate"`
	ReportReceivedDate       string `json:"reportReceivedDate"`
	Outcome     string `json:"outcome"`
	ReportReviewDate        string `json:"reportReviewDate"`
	ReviewedBy  string `json:"reviewedBy"`
	ReportMarkedAs  string `json:"reportMarkedAs"`
}

func (c *Client) AddAssuranceVisit(ctx Context, requestedDate string, userId int, deputyId int) error {
	var k FirmDetails
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(AssuranceVisit{
		RequestedDate:    requestedDate,
		RequestedById:     userId,
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/deputy/%d/assurance-visit", deputyId), &body)
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

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&k)
	return err
}
