package sirius

import (
	"encoding/json"
	"fmt"
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
	CurrentEcm  int
}

type Team struct {
	Members []TeamMember
}

func (c *Client) GetDeputyTeamMembers(ctx Context, defaultPATeam int, deputyDetails DeputyDetails) ([]TeamMember, error) {
	requestUrl := getRequestURL(deputyDetails, defaultPATeam)
	req, err := c.newRequest(ctx, http.MethodGet, requestUrl, nil)
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

	teamMembers := getTeamMembersByDeputyType(deputyDetails, resp)

	return teamMembers, nil
}

func getRequestURL(deputyDetails DeputyDetails, defaultPATeam int) string {
	if deputyDetails.DeputyType.Handle == "PRO" {
		return fmt.Sprintf("/api/v1/teams?type=%s", deputyDetails.DeputyType.Handle)
	} else {
		return "/api/v1/teams/" + strconv.Itoa(defaultPATeam)
	}
}

func getTeamMembersByDeputyType(deputyDetails DeputyDetails, resp *http.Response) []TeamMember {
	if deputyDetails.DeputyType.Handle == "PRO" {
		var v []apiTeam
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return []TeamMember{}
		}

		members := []TeamMember{}

		for _, k := range v {
			for _, m := range k.Members {
				members = append(members, TeamMember{
					ID:          m.ID,
					DisplayName: m.DisplayName,
					CurrentEcm:  deputyDetails.ExecutiveCaseManager.EcmId,
				})
			}
		}

		return members
	} else {
		var v apiTeam
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return []TeamMember{}
		}

		team := Team{}

		for _, m := range v.Members {
			team.Members = append(team.Members, TeamMember{
				ID:          m.ID,
				DisplayName: m.DisplayName,
			})
		}

		return team.Members
	}
}
