package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type addNoteRequest struct {
	Title       string `json:"name"`
	Note        string `json:"description"`
	UserId		int    `json:"createdById"`
	NoteType	string    `json:"noteType"`
}

func (c *Client) AddNote(ctx Context, title, note string, deputyId int) (int, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(addNoteRequest{
		Title:        title,
		Note:         note,
		UserId: 	  68,
		NoteType:     "PA_DEPUTY_NOTE_CREATED",
	})
	if err != nil {
		return 0, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/deputy/%d/create-note", deputyId), &body)

	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return 0, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusCreated {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return 0, ValidationError{Errors: v.ValidationErrors}
		}

		return 0, newStatusError(resp)
	}

	var v DeputyNote
	err = json.NewDecoder(resp.Body).Decode(&v)

	return v.DeputyId, err
}