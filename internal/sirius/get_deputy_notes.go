package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyNoteCollection []DeputyNote

type DeputyNote struct {
	ID              int    `json:"id"`
	DeputyId        int    `json:"personId"`
	UserId          int    `json:"userId"`
	UserDisplayName string `json:"userDisplayName"`
	UserEmail       string `json:"userEmail"`
	UserPhoneNumber string `json:"userPhoneNumber"`
	Type            string `json:"type"`
	NoteType        string `json:"noteType"`
	NoteText        string `json:"description"`
	Name            string `json:"name"`
	Timestamp       string `json:"createdTime"`
	Direction       string `json:"direction"`
}

func (c *Client) GetDeputyNotes(ctx Context, deputyId int) (DeputyNoteCollection, error) {
	var v DeputyNoteCollection

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(SupervisionAPIPath + "/v1/deputies/%d/notes", deputyId), nil)

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
