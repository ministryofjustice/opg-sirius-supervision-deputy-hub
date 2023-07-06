package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ContactDetails struct {
	ContactName      string `json:"contactName"`
	JobTitle         string `json:"jobTitle"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	OtherPhoneNumber string `json:"otherPhoneNumber"`
	Notes            string `json:"notes"`
	IsNamedDeputy    string `json:"isNamedDeputy"`
	IsMainContact    string `json:"isMainContact"`
}

func (c *Client) AddContactDetails(ctx Context, deputyId int, addContactForm ContactDetails) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(addContactForm)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/api/v1/deputies/%d/contacts", deputyId)

	req, err := c.newRequest(ctx, http.MethodPost, url, &body)

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

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
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

	return err
}
