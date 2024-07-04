package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type AssurancesList struct {
	Assurances []model.Assurance `json:"assurances"`
}

func (c *ApiClient) GetAssurances(ctx Context, deputyId int) ([]model.Assurance, error) {
	var k AssurancesList

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/assurances", deputyId), nil)

	if err != nil {
		return k.Assurances, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return k.Assurances, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return k.Assurances, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return k.Assurances, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&k); err != nil {
		return k.Assurances, err
	}

	return formatAssurances(k.Assurances, deputyId), nil
}

func formatAssurances(k []model.Assurance, deputyId int) []model.Assurance {
	var list []model.Assurance
	for _, s := range k {
		assurance := model.Assurance{
			Type:               s.Type,
			RequestedDate:      FormatDateTime(IsoDateTimeZone, s.RequestedDate, SiriusDate),
			Id:                 s.Id,
			RequestedBy:        s.RequestedBy,
			DeputyId:           deputyId,
			CommissionedDate:   FormatDateTime(IsoDateTimeZone, s.CommissionedDate, SiriusDate),
			ReportDueDate:      FormatDateTime(IsoDateTimeZone, s.ReportDueDate, SiriusDate),
			ReportReceivedDate: FormatDateTime(IsoDateTimeZone, s.ReportReceivedDate, SiriusDate),
			ReportReviewDate:   FormatDateTime(IsoDateTimeZone, s.ReportReviewDate, SiriusDate),
			VisitOutcome:       s.VisitOutcome,
			PdrOutcome:         s.PdrOutcome,
			Note:               s.Note,
			ReportMarkedAs:     s.ReportMarkedAs,
			VisitorAllocated:   s.VisitorAllocated,
			ReviewedBy:         s.ReviewedBy,
		}

		list = append(list, assurance)
	}
	return list
}
