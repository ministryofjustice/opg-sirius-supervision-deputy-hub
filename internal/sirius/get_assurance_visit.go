package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AssuranceVisit struct {
	Id                  int                 `json:"id"`
	AssuranceType       AssuranceTypes      `json:"assuranceType"`
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
	AssuranceVisitFormatted := formatAssuranceVisit(v)

	return AssuranceVisitFormatted, err
}

func formatAssuranceVisit(v AssuranceVisit) AssuranceVisit {
	updatedVisit := AssuranceVisit{
		AssuranceType:       v.AssuranceType,
		RequestedDate:       FormatDateAndTime(DateTimeFormat, v.RequestedDate, DateTimeDisplayFormat),
		Id:                  v.Id,
		RequestedBy:         v.RequestedBy,
		CommissionedDate:    FormatDateAndTime(DateTimeFormat, v.CommissionedDate, DateTimeDisplayFormat),
		ReportDueDate:       FormatDateAndTime(DateTimeFormat, v.ReportDueDate, DateTimeDisplayFormat),
		ReportReceivedDate:  FormatDateAndTime(DateTimeFormat, v.ReportReceivedDate, DateTimeDisplayFormat),
		ReportReviewDate:    FormatDateAndTime(DateTimeFormat, v.ReportReviewDate, DateTimeDisplayFormat),
		VisitOutcome:        v.VisitOutcome,
		VisitReportMarkedAs: v.VisitReportMarkedAs,
		VisitorAllocated:    v.VisitorAllocated,
		ReviewedBy:          v.ReviewedBy,
	}

	return updatedVisit
}
