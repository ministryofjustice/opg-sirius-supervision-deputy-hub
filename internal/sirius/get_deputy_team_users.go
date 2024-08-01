package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
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

func (c *ApiClient) GetDeputyTeamMembers(ctx Context, defaultPATeam int, deputyDetails DeputyDetails) ([]model.TeamMember, error) {

	requestUrl := getRequestURL(deputyDetails, defaultPATeam)
	req, err := c.newRequest(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return []model.TeamMember{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return []model.TeamMember{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return []model.TeamMember{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return []model.TeamMember{}, newStatusError(resp)
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

func getTeamMembersByDeputyType(deputyDetails DeputyDetails, resp *http.Response) []model.TeamMember {
	if deputyDetails.DeputyType.Handle == "PRO" {
		var v []apiTeam
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return []model.TeamMember{}
		}

		var members []model.TeamMember

		for _, k := range v {
			for _, m := range k.Members {
				members = append(members, model.TeamMember{
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
			return []model.TeamMember{}
		}

		team := model.Team{}

		for _, m := range v.Members {
			team.Members = append(team.Members, model.TeamMember{
				ID:          m.ID,
				DisplayName: m.DisplayName,
				CurrentEcm:  deputyDetails.ExecutiveCaseManager.EcmId,
			})
		}

		return team.Members
	}
}
