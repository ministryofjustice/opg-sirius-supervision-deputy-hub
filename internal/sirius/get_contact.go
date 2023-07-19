package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//type Contact struct {
//	Id                  int                 `json:"id"`
//	AssuranceType       AssuranceTypes      `json:"assuranceType"`
//	RequestedDate       string              `json:"requestedDate"`
//	RequestedBy         User                `json:"requestedBy"`
//	CommissionedDate    string              `json:"commissionedDate"`
//	ReportDueDate       string              `json:"reportDueDate"`
//	ReportReceivedDate  string              `json:"reportReceivedDate"`
//	VisitOutcome        VisitOutcomeTypes   `json:"assuranceVisitOutcome"`
//	PdrOutcome          PdrOutcomeTypes     `json:"pdrOutcome"`
//	ReportReviewDate    string              `json:"reportReviewDate"`
//	VisitReportMarkedAs VisitRagRatingTypes `json:"assuranceVisitReportMarkedAs"`
//	Note                string              `json:"note"`
//	VisitorAllocated    string              `json:"visitorAllocated"`
//	ReviewedBy          User                `json:"reviewedBy"`
//}

func (c *Client) GetContactById(ctx Context, deputyId int, contactId int) (Contact, error) {
	var contact Contact

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/contacts/%d", deputyId, contactId), nil)

	if err != nil {
		return contact, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return contact, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return contact, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return contact, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&contact)

	return contact, err
}
