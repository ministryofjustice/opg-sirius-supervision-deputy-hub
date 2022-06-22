package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AssuranceVisits struct {
	RequestedDate string `json:"requestedDate"`
	RequestedBy   User   `json:"requestedBy"`
	VisitId       int    `json:"id"`
	DeputyId      int
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
	AssuranceVisitsFormatted := editAssuranceVisits(k.AssuranceVisits, deputyId)

	return AssuranceVisitsFormatted, err
}

func editAssuranceVisits(k []AssuranceVisits, deputyId int) []AssuranceVisits {
	var list []AssuranceVisits
	for _, s := range k {
		event := AssuranceVisits{
			RequestedDate: FormatDateAndTime("2006-01-02T15:04:05+00:00", s.RequestedDate, "02/01/2006"),
			VisitId:       s.VisitId,
			RequestedBy:   s.RequestedBy,
			DeputyId:      deputyId,
		}

		list = append(list, event)
	}
	return list
}
