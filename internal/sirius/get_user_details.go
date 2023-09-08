package sirius

import (
	"encoding/json"
	"net/http"
	"strings"
)

type UserDetails struct {
	ID       int      `json:"id"`
	Roles    []string `json:"roles"`
	Username string   `json:"displayName"`
}

func (u UserDetails) IsFinanceManager() bool {
	for _, role := range u.Roles {
		if role == "Finance Manager" {
			return true
		}
	}

	return false
}

func (u UserDetails) IsSystemManager() bool {
	for _, role := range u.Roles {
		if role == "System Admin" {
			return true
		}
	}

	return false
}

func (u UserDetails) GetRoles() string {
	return strings.Join(u.Roles, ",")
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
