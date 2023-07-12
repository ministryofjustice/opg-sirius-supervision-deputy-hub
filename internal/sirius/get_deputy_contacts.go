package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	
)

type apiContact struct {
	ContactId        int    `json:"id"`
	Name             string `json:"name"`
	JobTitle         string `json:"jobTitle"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	OtherPhoneNumber string `json:"otherPhoneNumber"`
	ContactNotes     string `json:"contactNotes"`
	IsMainContact    bool   `json:"isMainContact"`
	IsNamedDeputy    bool   `json:"isNamedDeputy"`
}

type ApiContactList struct {
	Contacts      []apiContact `json:"contacts"`
	Pages         Page         `json:"pages"`
	Metadata      Metadata     `json:"metadata"`
	TotalContacts int          `json:"total"`
}

type ContactList struct {
	Contacts      DeputyContactsDetails
	Pages         Page
	TotalContacts int
	Metadata      Metadata
}

type DeputyContact struct {
	ContactId        int
	Name             string
	JobTitle         string
	Email            string
	PhoneNumber      string
	OtherPhoneNumber string
	ContactNotes     string
	IsMainContact    bool
	IsNamedDeputy    bool
}

type DeputyContactsDetails []DeputyContact

func (c *Client) GetDeputyContacts(ctx Context, deputyId, displayContactLimit, search int, deputyType, columnBeingSorted, sortOrder string) (ContactList, AriaSorting, error) {
	var contactList ContactList
	var apiContactList ApiContactList

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/contacts?&limit=%d&page=%d", deputyId, displayContactLimit, search), nil)

	if err != nil {
		return contactList, AriaSorting{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return contactList, AriaSorting{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return contactList, AriaSorting{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return contactList, AriaSorting{}, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&apiContactList); err != nil {
		return contactList, AriaSorting{}, err
	}

	var contacts DeputyContactsDetails
	for _, t := range apiContactList.Contacts {

			var contact = DeputyContact{
				ContactId:        t.ContactId,
				Name:             t.Name,
				JobTitle:         t.JobTitle,
				Email:            t.Email,
				PhoneNumber:      t.PhoneNumber,
				OtherPhoneNumber: t.OtherPhoneNumber,
				ContactNotes:     t.ContactNotes,
				IsMainContact:    t.IsMainContact,
				IsNamedDeputy:    t.IsNamedDeputy,
			}
			
			contacts = append(contacts, contact)
	}
	contactList.Contacts = contacts

	var aria AriaSorting
	aria.SurnameAriaSort = changeSortButtonDirection(sortOrder, columnBeingSorted, "surname")
	aria.ReportDueAriaSort = changeSortButtonDirection(sortOrder, columnBeingSorted, "reportdue")
	aria.CRECAriaSort = changeSortButtonDirection(sortOrder, columnBeingSorted, "crec")

	return contactList, aria, err
}