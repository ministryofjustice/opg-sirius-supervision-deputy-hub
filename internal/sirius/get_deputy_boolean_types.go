package sirius

import (
	"encoding/json"
	"net/http"
)

type DeputyBooleanTypes struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *Client) GetDeputyBooleanTypes(ctx Context) ([]DeputyBooleanTypes, error) {
	var v []DeputyBooleanTypes

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/reference-data/deputyBooleanType", nil)
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
