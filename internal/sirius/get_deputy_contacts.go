package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiContact struct {
	Id                            int    `json:"id"`
	Name                          string `json:"name"`
	JobTitle                      string `json:"jobTitle"`
	Email                         string `json:"email"`
	PhoneNumber                   string `json:"phoneNumber"`
	OtherPhoneNumber              string `json:"otherPhoneNumber"`
	Notes                         string `json:"notes"`
	IsMainContact                 bool   `json:"isMainContact"`
	IsNamedDeputy                 bool   `json:"isNamedDeputy"`
	IsMonthlySpreadsheetRecipient bool   `json:"isMonthlySpreadsheetRecipient"`
}

type DeputyContact struct {
	Id                            int
	Name                          string
	JobTitle                      string
	Email                         string
	PhoneNumber                   string
	OtherPhoneNumber              string
	Notes                         string
	IsMainContact                 bool
	IsNamedDeputy                 bool
	IsMonthlySpreadsheetRecipient bool
}

type ContactList []DeputyContact

func (c *Client) GetDeputyContacts(ctx Context, deputyId int) (ContactList, error) {
	var contactList ContactList
	var apiContacts []ApiContact

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/contacts", deputyId), nil)

	if err != nil {
		return contactList, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return contactList, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return contactList, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return contactList, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&apiContacts); err != nil {
		return contactList, err
	}

	for _, t := range apiContacts {
		contactList = append(contactList, DeputyContact(t))
	}

	return contactList, err
}
