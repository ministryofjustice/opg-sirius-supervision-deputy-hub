package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetAssuranceVisits(ctx Context, deputyId int) (AssuranceVisit, error) {
	var k AssuranceVisit

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/assurance-visit", deputyId), nil)

	if err != nil {
		return k, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return k, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return k, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return k, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&k)
	return k, err
}
