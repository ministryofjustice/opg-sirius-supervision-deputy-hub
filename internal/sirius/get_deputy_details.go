package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExecutiveCaseManager struct {
	EcmId   int    `json:"id"`
	EcmName string `json:"displayName"`
}

type ExecutiveCaseManagerOutgoing struct {
	EcmId   int
}

type DeputyDetails struct {
	ID                               int                  `json:"id"`
	DeputyCasrecId                   int                  `json:"deputyCasrecId"`
	DeputyNumber                     int                  `json:"deputyNumber"`
	OrganisationName                 string               `json:"organisationName"`
	OrganisationTeamOrDepartmentName string               `json:"organisationTeamOrDepartmentName"`
	Email                            string               `json:"email"`
	PhoneNumber                      string               `json:"phoneNumber"`
	AddressLine1                     string               `json:"addressLine1"`
	AddressLine2                     string               `json:"addressLine2"`
	AddressLine3                     string               `json:"addressLine3"`
	Town                             string               `json:"town"`
	County                           string               `json:"county"`
	Postcode                         string               `json:"postcode"`
	ExecutiveCaseManager             ExecutiveCaseManager `json:"executiveCaseManager"`
}

func (c *Client) GetDeputyDetails(ctx Context, defaultPATeam int, deputyId int) (DeputyDetails, error) {
	var v DeputyDetails

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d", deputyId), nil)
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

	if v.ExecutiveCaseManager.EcmId == 0 {
		v.ExecutiveCaseManager.EcmId = defaultPATeam
		v.ExecutiveCaseManager.EcmName = "Public Authority Deputy Team"
	}
	return v, err
}
