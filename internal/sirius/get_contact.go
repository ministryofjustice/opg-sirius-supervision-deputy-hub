package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Contact struct {
	ContactName                   string `json:"name"`
	JobTitle                      string `json:"jobTitle"`
	Email                         string `json:"email"`
	PhoneNumber                   string `json:"phoneNumber"`
	OtherPhoneNumber              string `json:"otherPhoneNumber"`
	ContactNotes                  string `json:"notes"`
	IsNamedDeputy                 bool   `json:"isNamedDeputy"`
	IsMainContact                 bool   `json:"isMainContact"`
	IsMonthlySpreadsheetRecipient bool   `json:"isMonthlySpreadsheetRecipient"`
}

func (c *Client) GetContactById(ctx Context, deputyId int, contactId int) (Contact, error) {
	var contact Contact

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/contacts/%d", deputyId, contactId), nil)

	if err != nil {
		return contact, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return contact, err
	}
	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return contact, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return contact, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&contact)

	return contact, err
}
