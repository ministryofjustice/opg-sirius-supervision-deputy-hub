package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type AssuranceVisits struct {
	AssuranceType       AssuranceTypes      `json:"assuranceType"`
	RequestedDate       string              `json:"requestedDate"`
	RequestedBy         model.User          `json:"requestedBy"`
	VisitId             int                 `json:"id"`
	CommissionedDate    string              `json:"commissionedDate"`
	ReportDueDate       string              `json:"reportDueDate"`
	ReportReceivedDate  string              `json:"reportReceivedDate"`
	VisitOutcome        VisitOutcomeTypes   `json:"assuranceVisitOutcome"`
	PdrOutcome          PdrOutcomeTypes     `json:"pdrOutcome"`
	Note                string              `json:"note"`
	ReportReviewDate    string              `json:"reportReviewDate"`
	VisitReportMarkedAs VisitRagRatingTypes `json:"assuranceVisitReportMarkedAs"`
	VisitorAllocated    string              `json:"visitorAllocated"`
	ReviewedBy          model.User          `json:"reviewedBy"`
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
			RequestedDate:       FormatDateTime(IsoDateTimeZone, s.RequestedDate, SiriusDate),
			VisitId:             s.VisitId,
			RequestedBy:         s.RequestedBy,
			DeputyId:            deputyId,
			CommissionedDate:    FormatDateTime(IsoDateTimeZone, s.CommissionedDate, SiriusDate),
			ReportDueDate:       FormatDateTime(IsoDateTimeZone, s.ReportDueDate, SiriusDate),
			ReportReceivedDate:  FormatDateTime(IsoDateTimeZone, s.ReportReceivedDate, SiriusDate),
			ReportReviewDate:    FormatDateTime(IsoDateTimeZone, s.ReportReviewDate, SiriusDate),
			VisitOutcome:        s.VisitOutcome,
			PdrOutcome:          s.PdrOutcome,
			Note:                s.Note,
			VisitReportMarkedAs: s.VisitReportMarkedAs,
			VisitorAllocated:    s.VisitorAllocated,
			ReviewedBy:          s.ReviewedBy,
		}

		list = append(list, event)
	}
	return list
}
