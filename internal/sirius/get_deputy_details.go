package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyDetails struct {
	ID               int    `json:"id"`
	DeputyCasrecId   int `json:"deputyCasrecId"`
	OrganisationName string `json:"organisationName"`
}

func (c *Client) GetDeputyDetails(ctx Context, deputyId int) (DeputyDetails, error) {
	var v DeputyDetails

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d", deputyId), nil)
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
	return v, err
}
