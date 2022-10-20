package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AssuranceVisitDetails struct {
	CommissionedDate    string `json:"commissionedDate"`
	ReportDueDate       string `json:"reportDueDate"`
	ReportReceivedDate  string `json:"reportReceivedDate"`
	VisitOutcome        string `json:"assuranceVisitOutcome"`
	PdrOutcome          string `json:"pdrOutcome"`
	ReportReviewDate    string `json:"reportReviewDate"`
	VisitReportMarkedAs string `json:"assuranceVisitReportMarkedAs"`
	VisitorAllocated    string `json:"visitorAllocated"`
	ReviewedBy          int    `json:"reviewedBy"`
}

func (c *Client) UpdateAssuranceVisit(ctx Context, manageAssuranceVisitForm AssuranceVisitDetails, deputyId, visitId int) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(AssuranceVisitDetails{
		CommissionedDate:    manageAssuranceVisitForm.CommissionedDate,
		ReportDueDate:       manageAssuranceVisitForm.ReportDueDate,
		ReportReceivedDate:  manageAssuranceVisitForm.ReportReceivedDate,
		VisitOutcome:        manageAssuranceVisitForm.VisitOutcome,
		PdrOutcome:          manageAssuranceVisitForm.PdrOutcome,
		ReportReviewDate:    manageAssuranceVisitForm.ReportReviewDate,
		VisitReportMarkedAs: manageAssuranceVisitForm.VisitReportMarkedAs,
		VisitorAllocated:    manageAssuranceVisitForm.VisitorAllocated,
		ReviewedBy:          manageAssuranceVisitForm.ReviewedBy,
	})

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/deputies/%d/assurance-visits/%d", deputyId, visitId), &body)

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
		return newStatusError(resp)
	}

	return nil
}
