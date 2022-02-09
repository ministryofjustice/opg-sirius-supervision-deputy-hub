package sirius

import (
	"encoding/json"
	"net/http"
)

type DeputyReportSystemTypes struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *Client) GetDeputyReportSystemTypes(ctx Context) ([]DeputyReportSystemTypes, error) {
	var v []DeputyReportSystemTypes

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/reference-data/deputyReportSystem", nil)
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
