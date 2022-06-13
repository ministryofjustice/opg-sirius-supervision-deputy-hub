package sirius

import (
	"encoding/json"
	"net/http"
)

type UserDetails struct {
	ID    int      `json:"id"`
	Roles []string `json:"roles"`
	Username string `json:"displayName"`
}

func (d UserDetails) IsFinanceManager() bool {
	for _, role := range d.Roles {
		if role == "Finance Manager" {
			return true
		}
	}

	return false
}

func (c *Client) GetUserDetails(ctx Context) (UserDetails, error) {
	var v UserDetails

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/users/current", nil)
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
