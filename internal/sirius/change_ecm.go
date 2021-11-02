package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ChangeECM(ctx Context, changeDeputyECMForm, deputyDetails DeputyDetails) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(editDeputyDetails{
		ID:                               changeDeputyECMForm.ID,
		OrganisationName:                 deputyDetails.OrganisationName,
		OrganisationTeamOrDepartmentName: changeDeputyECMForm.OrganisationTeamOrDepartmentName,
		Email:                            deputyDetails.Email,
		PhoneNumber:                      deputyDetails.PhoneNumber,
		AddressLine1:                     deputyDetails.AddressLine1,
		AddressLine2:                     deputyDetails.AddressLine2,
		AddressLine3:                     deputyDetails.AddressLine3,
		Town:                             deputyDetails.Town,
		County:                           deputyDetails.County,
		Postcode:                         deputyDetails.Postcode,
		DeputyType: Deputy{Handle: "PA",
			Label: "Public Authority"},
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/deputies/%d", changeDeputyECMForm.ID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
