package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Deputy struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type editDeputyDetails struct {
	ID                               int    `json:"id"`
	OrganisationName                 string `json:"organisationName"`
	OrganisationTeamOrDepartmentName string `json:"organisationTeamOrDepartmentName"`
	Email                            string `json:"email"`
	PhoneNumber                      string `json:"workPhoneNumber"`
	AddressLine1                     string `json:"addressLine1"`
	AddressLine2                     string `json:"addressLine2"`
	AddressLine3                     string `json:"addressLine3"`
	Town                             string `json:"town"`
	County                           string `json:"county"`
	Postcode                         string `json:"postcode"`
	Salutation                       string `json:"salutation"`
	Firstname                        string `json:"firstname"`
	OtherNames                       string `json:"otherNames"`
	Surname                          string `json:"surname"`
	Dob                              string `json:"dob"`
	PreviousNames                    string `json:"previousNames"`
	DeputyType                       Deputy `json:"deputyType"`
}

func (c *Client) EditDeputyDetails(ctx Context, editDeputyDetailForm DeputyDetails) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(editDeputyDetails{
		ID:                               editDeputyDetailForm.ID,
		OrganisationName:                 editDeputyDetailForm.OrganisationName,
		OrganisationTeamOrDepartmentName: editDeputyDetailForm.OrganisationTeamOrDepartmentName,
		Email:                            editDeputyDetailForm.Email,
		PhoneNumber:                      editDeputyDetailForm.PhoneNumber,
		AddressLine1:                     editDeputyDetailForm.AddressLine1,
		AddressLine2:                     editDeputyDetailForm.AddressLine2,
		AddressLine3:                     editDeputyDetailForm.AddressLine3,
		Town:                             editDeputyDetailForm.Town,
		County:                           editDeputyDetailForm.County,
		Postcode:                         editDeputyDetailForm.Postcode,
		DeputyType: Deputy{Handle: "PA",
			Label: "Public Authority"},
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d", editDeputyDetailForm.ID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
