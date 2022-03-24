package sirius

import (
	"encoding/json"
	"net/http"
)

type FirmForList struct {
	Id         int    `json:"id"`
	FirmName   string `json:"firmName"`
	FirmNumber int    `json:"firmNumber"`
}

func (c *Client) GetFirms(ctx Context) ([]FirmForList, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/firms", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v []FirmForList
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, err
}
