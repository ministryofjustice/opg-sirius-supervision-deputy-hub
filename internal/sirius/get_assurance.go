package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

func (c *ApiClient) GetAssuranceById(ctx Context, deputyId int, visitId int) (model.Assurance, error) {
	var v model.Assurance

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/assurances/%d", deputyId, visitId), nil)

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

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}

	return formatAssurance(v), nil
}

func formatAssurance(v model.Assurance) model.Assurance {
	return model.Assurance{
		Type:               v.Type,
		RequestedDate:      FormatDateTime(IsoDateTimeZone, v.RequestedDate, IsoDate),
		Id:                 v.Id,
		RequestedBy:        v.RequestedBy,
		CommissionedDate:   FormatDateTime(IsoDateTimeZone, v.CommissionedDate, IsoDate),
		ReportDueDate:      FormatDateTime(IsoDateTimeZone, v.ReportDueDate, IsoDate),
		ReportReceivedDate: FormatDateTime(IsoDateTimeZone, v.ReportReceivedDate, IsoDate),
		ReportReviewDate:   FormatDateTime(IsoDateTimeZone, v.ReportReviewDate, IsoDate),
		VisitOutcome:       v.VisitOutcome,
		PdrOutcome:         v.PdrOutcome,
		ReportMarkedAs:     v.ReportMarkedAs,
		Note:               v.Note,
		VisitorAllocated:   v.VisitorAllocated,
		ReviewedBy:         v.ReviewedBy,
	}
}
