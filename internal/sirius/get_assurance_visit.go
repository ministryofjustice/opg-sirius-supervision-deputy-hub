package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type AssuranceVisit struct {
	Id                  int                 `json:"id"`
	AssuranceType       AssuranceTypes      `json:"assuranceType"`
	RequestedDate       string              `json:"requestedDate"`
	RequestedBy         model.User          `json:"requestedBy"`
	CommissionedDate    string              `json:"commissionedDate"`
	ReportDueDate       string              `json:"reportDueDate"`
	ReportReceivedDate  string              `json:"reportReceivedDate"`
	VisitOutcome        VisitOutcomeTypes   `json:"assuranceVisitOutcome"`
	PdrOutcome          PdrOutcomeTypes     `json:"pdrOutcome"`
	ReportReviewDate    string              `json:"reportReviewDate"`
	VisitReportMarkedAs VisitRagRatingTypes `json:"assuranceVisitReportMarkedAs"`
	Note                string              `json:"note"`
	VisitorAllocated    string              `json:"visitorAllocated"`
	ReviewedBy          model.User          `json:"reviewedBy"`
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
		RequestedDate:       FormatDateTime(IsoDateTimeZone, v.RequestedDate, IsoDate),
		Id:                  v.Id,
		RequestedBy:         v.RequestedBy,
		CommissionedDate:    FormatDateTime(IsoDateTimeZone, v.CommissionedDate, IsoDate),
		ReportDueDate:       FormatDateTime(IsoDateTimeZone, v.ReportDueDate, IsoDate),
		ReportReceivedDate:  FormatDateTime(IsoDateTimeZone, v.ReportReceivedDate, IsoDate),
		ReportReviewDate:    FormatDateTime(IsoDateTimeZone, v.ReportReviewDate, IsoDate),
		VisitOutcome:        v.VisitOutcome,
		PdrOutcome:          v.PdrOutcome,
		VisitReportMarkedAs: v.VisitReportMarkedAs,
		Note:                v.Note,
		VisitorAllocated:    v.VisitorAllocated,
		ReviewedBy:          v.ReviewedBy,
	}

	return updatedVisit
}
