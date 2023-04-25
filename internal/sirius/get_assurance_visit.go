package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AssuranceVisit struct {
	Id                  int                 `json:"id"`
	AssuranceType       AssuranceTypes      `json:"assuranceType"`
	RequestedDate       time.Time           `json:"requestedDate"`
	RequestedBy         User                `json:"requestedBy"`
	CommissionedDate    time.Time           `json:"commissionedDate"`
	ReportDueDate       time.Time           `json:"reportDueDate"`
	ReportReceivedDate  time.Time           `json:"reportReceivedDate"`
	VisitOutcome        VisitOutcomeTypes   `json:"assuranceVisitOutcome"`
	PdrOutcome          PdrOutcomeTypes     `json:"pdrOutcome"`
	ReportReviewDate    time.Time           `json:"reportReviewDate"`
	VisitReportMarkedAs VisitRagRatingTypes `json:"assuranceVisitReportMarkedAs"`
	Note                string              `json:"note"`
	VisitorAllocated    string              `json:"visitorAllocated"`
	ReviewedBy          User                `json:"reviewedBy"`
	NullDateForFrontEnd time.Time
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
		RequestedDate:       v.RequestedDate,
		Id:                  v.Id,
		RequestedBy:         v.RequestedBy,
		CommissionedDate:    v.CommissionedDate,
		ReportDueDate:       v.ReportDueDate,
		ReportReceivedDate:  v.ReportReceivedDate,
		ReportReviewDate:    v.ReportReviewDate,
		VisitOutcome:        v.VisitOutcome,
		PdrOutcome:          v.PdrOutcome,
		VisitReportMarkedAs: v.VisitReportMarkedAs,
		Note:                v.Note,
		VisitorAllocated:    v.VisitorAllocated,
		ReviewedBy:          v.ReviewedBy,
		NullDateForFrontEnd: GetNullDate(),
	}

	return updatedVisit
}
