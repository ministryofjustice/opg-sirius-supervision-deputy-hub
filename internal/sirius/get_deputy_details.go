package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type executiveCaseManager struct {
	EcmId   int    `json:"id"`
	EcmName string `json:"displayName"`
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
	ExecutiveCaseManager             executiveCaseManager `json:"executiveCaseManager"`
}

func (c *Client) GetDeputyDetails(ctx Context, defaultPATeam string, deputyId int) (DeputyDetails, error) {
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

	return v, err
}
