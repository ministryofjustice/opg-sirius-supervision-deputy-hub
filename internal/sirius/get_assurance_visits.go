package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AssuranceVisits struct {
	AssuranceType       AssuranceTypes      `json:"assuranceType"`
	RequestedDate       time.Time           `json:"requestedDate"`
	RequestedBy         User                `json:"requestedBy"`
	VisitId             int                 `json:"id"`
	CommissionedDate    time.Time           `json:"commissionedDate"`
	ReportDueDate       time.Time           `json:"reportDueDate"`
	ReportReceivedDate  time.Time           `json:"reportReceivedDate"`
	VisitOutcome        VisitOutcomeTypes   `json:"assuranceVisitOutcome"`
	PdrOutcome          PdrOutcomeTypes     `json:"pdrOutcome"`
	Note                string              `json:"note"`
	ReportReviewDate    time.Time           `json:"reportReviewDate"`
	VisitReportMarkedAs VisitRagRatingTypes `json:"assuranceVisitReportMarkedAs"`
	VisitorAllocated    string              `json:"visitorAllocated"`
	ReviewedBy          User                `json:"reviewedBy"`
	DeputyId            int
	NullDateForFrontEnd time.Time
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
	nullDate, _ := time.Parse("2006-01-02T15:04:05+00:00", "0001-01-01 00:00:00 +0000 UTC")

	for _, s := range k {
		event := AssuranceVisits{
			AssuranceType:       s.AssuranceType,
			RequestedDate:       s.RequestedDate,
			VisitId:             s.VisitId,
			RequestedBy:         s.RequestedBy,
			DeputyId:            deputyId,
			CommissionedDate:    s.CommissionedDate,
			ReportDueDate:       s.ReportDueDate,
			ReportReceivedDate:  s.ReportReceivedDate,
			ReportReviewDate:    s.ReportReviewDate,
			VisitOutcome:        s.VisitOutcome,
			PdrOutcome:          s.PdrOutcome,
			Note:                s.Note,
			VisitReportMarkedAs: s.VisitReportMarkedAs,
			VisitorAllocated:    s.VisitorAllocated,
			ReviewedBy:          s.ReviewedBy,
			NullDateForFrontEnd: nullDate,
		}

		list = append(list, event)
	}
	return list
}
