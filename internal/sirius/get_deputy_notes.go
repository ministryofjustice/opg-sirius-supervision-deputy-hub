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

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputy/%d/notes", deputyId), nil)

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

	DeputyNotes := EditDeputyNotes(v)

	return DeputyNotes, err
}

func ReformatTimestampDeputyNote(s DeputyNote) string {
	return s.Timestamp
}

func EditDeputyNotes(v DeputyNoteCollection) DeputyNoteCollection {
	var list DeputyNoteCollection
	for _, s := range v {
		note := DeputyNote{
			DeputyId:        s.DeputyId,
			UserId:          s.UserId,
			UserDisplayName: s.UserDisplayName,
			UserEmail:       s.UserEmail,
			UserPhoneNumber: s.UserPhoneNumber,
			ID:              s.ID,
			Type:            s.Type,
			NoteType:        s.NoteType,
			NoteText:     s.NoteText,
			Name:            s.Name,
			Timestamp:       s.Timestamp,
			Direction:       s.Direction,
		}

		list = append(list, note)
	}
	return list
}

