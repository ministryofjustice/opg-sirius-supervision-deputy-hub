package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiContact struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	JobTitle         string `json:"jobTitle"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	OtherPhoneNumber string `json:"otherPhoneNumber"`
	Notes            string `json:"notes"`
	IsMainContact    bool   `json:"isMainContact"`
	IsNamedDeputy    bool   `json:"isNamedDeputy"`
}

type ContactList struct {
	Contacts      DeputyContactsDetails
	Pages         Page
	TotalContacts int
	Metadata      Metadata
	DeputyId      int
}

type DeputyContact struct {
	Id               int
	Name             string
	JobTitle         string
	Email            string
	PhoneNumber      string
	OtherPhoneNumber string
	Notes            string
	IsMainContact    bool
	IsNamedDeputy    bool
}

type DeputyContactsDetails []DeputyContact

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

	var contacts DeputyContactsDetails
	for _, t := range apiContacts {
		contacts = append(contacts, DeputyContact(t))
	}
	contactList.Contacts = contacts
	contactList.DeputyId = deputyId

	return contactList, err
}
