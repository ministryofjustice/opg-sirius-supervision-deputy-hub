package sirius

import (
"encoding/json"
"net/http"
"strconv"
)

type apiTeam struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Members     []struct {
		ID          int    `json:"id"`
		DisplayName string `json:"displayName"`
		Email       string `json:"email"`
	} `json:"members"`
	TeamType *struct {
		Handle string `json:"handle"`
		Label  string `json:"label"`
	} `json:"teamType"`
}

type TeamMember struct {
	ID          int
	DisplayName string
	Email       string
}

type Team struct {
	Members     []TeamMember
}


func (c *Client) GetPaDeputyTeamMembers(ctx Context, defaultPATeam int) ([]TeamMember, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams/"+strconv.Itoa(defaultPATeam), nil)
	if err != nil {
		return []TeamMember{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return []TeamMember{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return []TeamMember{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return []TeamMember{}, newStatusError(resp)
	}

	var v apiTeam
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return []TeamMember{}, err
	}

	team := Team{}

	for _, m := range v.Members {
		team.Members = append(team.Members, TeamMember{
			ID:          m.ID,
			DisplayName: m.DisplayName,
		})
	}

	return team.Members, nil
}