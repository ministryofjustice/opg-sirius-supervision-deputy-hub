package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyContactDetails struct {
	OrganisationName string `json:"organisationName"`
	DeputyFirstName  string `json:"firstname"`
	DeputySurname    string `json:"surname"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"workPhoneNumber"`
	AddressLine1     string `json:"addressLine1"`
	AddressLine2     string `json:"addressLine2"`
	AddressLine3     string `json:"addressLine3"`
	Town             string `json:"town"`
	County           string `json:"county"`
	Country          string `json:"country"`
	Postcode         string `json:"postcode"`
	DeputySubType    string `json:"deputySubType"`
}

func (c *Client) UpdateDeputyContactDetails(ctx Context, deputyId int, deputyDetails DeputyContactDetails) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(deputyDetails)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/deputies/%d/contact-details", deputyId)

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

	if resp.StatusCode == http.StatusBadRequest {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
