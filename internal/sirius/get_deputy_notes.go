package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyNoteList struct {
	Limit	int 	`json:"limit"`
	Pages	Pages 	`json:"pages"`
	Total 	int 	`json:"total"`
	DeputyNotes []DeputyNote	`json:"notes"`
}

type Pages struct {
	Current int	`json:"current"`
	Total	int `json:"total"`
}

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

func (c *Client) GetDeputyNotes(ctx Context, deputyId int) (DeputyNoteList, error) {
	var v DeputyNoteList

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

