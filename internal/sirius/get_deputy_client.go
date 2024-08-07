package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClientWithOrderDeputy struct {
	ClientId  int    `json:"id"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	CourtRef  string `json:"caseRecNumber"`
	Cases     []struct {
		Deputies []struct {
			Deputy struct {
				Id int `json:"id"`
			} `json:"deputy"`
		} `json:"deputies"`
	} `json:"cases"`
}

func (c *Client) GetDeputyClient(ctx Context, deputyId int, caseRecNumber string) (ClientWithOrderDeputy, error) {
	var v ClientWithOrderDeputy

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/client/%s", deputyId, caseRecNumber), nil)

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

	return v, err
}
