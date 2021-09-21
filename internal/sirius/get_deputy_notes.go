package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyNoteCollection []DeputyNote

type DeputyNote struct {
	ID               int    `json:"id"`
	DeputyCasrecId   int `json:"deputyCasrecId"`
	OrganisationName string `json:"organisationName"`
	Email	string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	AddressLine1	string `json:"addressLine1"`
	AddressLine2	string `json:"addressLine2"`
	AddressLine3	string `json:"addressLine3"`
	Town string `json:"town"`
	County string `json:"county"`
	Postcode string `json:"postcode"`
}

func (c *Client) GetDeputyNotes(ctx Context, deputyId int) (DeputyNoteCollection, error) {
	var v DeputyNoteCollection

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/clients/%d/notes", deputyId), nil)

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