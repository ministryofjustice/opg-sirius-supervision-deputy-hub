package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AssuranceVisit struct {
	Id                  int                 `json:"id"`
	RequestedDate       string              `json:"requestedDate"`
	RequestedBy         User                `json:"requestedBy"`
	CommissionedDate    string              `json:"commissionedDate"`
	ReportDueDate       string              `json:"reportDueDate"`
	ReportReceivedDate  string              `json:"reportReceivedDate"`
	VisitOutcome        VisitOutcomeTypes   `json:"assuranceVisitOutcome"`
	ReportReviewDate    string              `json:"reportReviewDate"`
	VisitReportMarkedAs VisitRagRatingTypes `json:"assuranceVisitReportMarkedAs"`
	VisitorAllocated    string              `json:"visitorAllocated"`
	ReviewedBy          User                `json:"reviewedBy"`
}

func (c *Client) GetAssuranceVisitById(ctx Context, deputyId int, visitId int) (AssuranceVisit, error) {
	var v AssuranceVisit

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/assurance-visits/%d", deputyId, visitId), nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	AssuranceVisitFormatted := editAssuranceVisit(v)

	return AssuranceVisitFormatted, err
}

func editAssuranceVisit(v AssuranceVisit) AssuranceVisit {
	updatedVisit := AssuranceVisit{
		RequestedDate:       FormatDateAndTime("2006-01-02T15:04:05+00:00", v.RequestedDate, "2006-01-02"),
		Id:                  v.Id,
		RequestedBy:         v.RequestedBy,
		CommissionedDate:    FormatDateAndTime("2006-01-02T15:04:05+00:00", v.CommissionedDate, "2006-01-02"),
		ReportDueDate:       FormatDateAndTime("2006-01-02T15:04:05+00:00", v.ReportDueDate, "2006-01-02"),
		ReportReceivedDate:  FormatDateAndTime("2006-01-02T15:04:05+00:00", v.ReportReceivedDate, "2006-01-02"),
		ReportReviewDate:    FormatDateAndTime("2006-01-02T15:04:05+00:00", v.ReportReviewDate, "2006-01-02"),
		VisitOutcome:        v.VisitOutcome,
		VisitReportMarkedAs: v.VisitReportMarkedAs,
		VisitorAllocated:    v.VisitorAllocated,
		ReviewedBy:          v.ReviewedBy,
	}

	return updatedVisit
}
