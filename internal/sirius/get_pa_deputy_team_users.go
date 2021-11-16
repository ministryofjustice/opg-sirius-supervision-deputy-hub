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
	ID          int
	DisplayName string
	Members     []TeamMember
	Type        string
	TypeLabel   string
	Email       string
	PhoneNumber string
}


func (c *Client) GetPaDeputyTeamMembers(ctx Context, defaultPATeam int) (Team, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams/"+strconv.Itoa(defaultPATeam), nil)
	if err != nil {
		return Team{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Team{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return Team{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return Team{}, newStatusError(resp)
	}

	var v apiTeam
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return Team{}, err
	}

	team := Team{
		ID:          v.ID,
		DisplayName: v.DisplayName,
		Type:        "",
		Email:       v.Email,
		PhoneNumber: v.PhoneNumber,
	}

	for _, m := range v.Members {
		team.Members = append(team.Members, TeamMember{
			ID:          m.ID,
			DisplayName: m.DisplayName,
			Email:       m.Email,
		})
	}

	if v.TeamType != nil {
		team.Type = v.TeamType.Handle
	}

	return team, nil
}