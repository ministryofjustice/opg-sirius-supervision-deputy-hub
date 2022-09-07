package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AssuranceVisits struct {
	AssuranceType       AssuranceTypes      `json:"assuranceType"`
	RequestedDate       string              `json:"requestedDate"`
	RequestedBy         User                `json:"requestedBy"`
	VisitId             int                 `json:"id"`
	CommissionedDate    string              `json:"commissionedDate"`
	ReportDueDate       string              `json:"reportDueDate"`
	ReportReceivedDate  string              `json:"reportReceivedDate"`
	VisitOutcome        VisitOutcomeTypes   `json:"assuranceVisitOutcome"`
	ReportReviewDate    string              `json:"reportReviewDate"`
	VisitReportMarkedAs VisitRagRatingTypes `json:"assuranceVisitReportMarkedAs"`
	VisitorAllocated    string              `json:"visitorAllocated"`
	ReviewedBy          User                `json:"reviewedBy"`
	DeputyId            int
}

type AssuranceVisitsList struct {
	AssuranceVisits []AssuranceVisits `json:"assuranceVisits"`
}

func (c *Client) GetAssuranceVisits(ctx Context, deputyId int) ([]AssuranceVisits, error) {
	var k AssuranceVisitsList

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/assurance-visits", deputyId), nil)

	if err != nil {
		return k.AssuranceVisits, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return k.AssuranceVisits, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return k.AssuranceVisits, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return k.AssuranceVisits, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&k)
	AssuranceVisitsFormatted := formatAssuranceVisits(k.AssuranceVisits, deputyId)

	return AssuranceVisitsFormatted, err
}

func formatAssuranceVisits(k []AssuranceVisits, deputyId int) []AssuranceVisits {
	var list []AssuranceVisits
	for _, s := range k {
		event := AssuranceVisits{
			AssuranceType:       s.AssuranceType,
			RequestedDate:       FormatDateAndTime("2006-01-02T15:04:05+00:00", s.RequestedDate, "02/01/2006"),
			VisitId:             s.VisitId,
			RequestedBy:         s.RequestedBy,
			DeputyId:            deputyId,
			CommissionedDate:    FormatDateAndTime("2006-01-02T15:04:05+00:00", s.CommissionedDate, "02/01/2006"),
			ReportDueDate:       FormatDateAndTime("2006-01-02T15:04:05+00:00", s.ReportDueDate, "02/01/2006"),
			ReportReceivedDate:  FormatDateAndTime("2006-01-02T15:04:05+00:00", s.ReportReceivedDate, "02/01/2006"),
			ReportReviewDate:    FormatDateAndTime("2006-01-02T15:04:05+00:00", s.ReportReviewDate, "02/01/2006"),
			VisitOutcome:        s.VisitOutcome,
			VisitReportMarkedAs: s.VisitReportMarkedAs,
			VisitorAllocated:    s.VisitorAllocated,
			ReviewedBy:          s.ReviewedBy,
		}

		list = append(list, event)
	}
	return list
}
