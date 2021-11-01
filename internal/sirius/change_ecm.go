package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ChangeECM(ctx Context, editDeputyDetailForm DeputyDetails) error {
	var body bytes.Buffer
	// can keep org team and dept name but need to autofill rest of data
	err := json.NewEncoder(&body).Encode(editDeputyDetails{
		// ID:                               editDeputyDetailForm.ID,
		// OrganisationName:                 editDeputyDetailForm.OrganisationName,
		OrganisationTeamOrDepartmentName: editDeputyDetailForm.OrganisationTeamOrDepartmentName,
		// Email:                            editDeputyDetailForm.Email,
		// PhoneNumber:                      editDeputyDetailForm.PhoneNumber,
		// AddressLine1:                     editDeputyDetailForm.AddressLine1,
		// AddressLine2:                     editDeputyDetailForm.AddressLine2,
		// AddressLine3:                     editDeputyDetailForm.AddressLine3,
		// Town:                             editDeputyDetailForm.Town,
		// County:                           editDeputyDetailForm.County,
		// Postcode:                         editDeputyDetailForm.Postcode,
		DeputyType: Deputy{Handle: "PA",
			Label: "Public Authority"},
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/deputies/%d", editDeputyDetailForm.ID)

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
