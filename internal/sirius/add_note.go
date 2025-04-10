package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type addNoteRequest struct {
	Title    string `json:"name"`
	Note     string `json:"description"`
	UserId   int    `json:"createdById"`
	NoteType string `json:"noteType"`
}

func (c *Client) AddNote(ctx Context, title, note string, deputyId, userId int, deputyType string) error {
	var noteType = getNoteType(deputyType)
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(addNoteRequest{
		Title:    title,
		Note:     note,
		UserId:   userId,
		NoteType: noteType,
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(SupervisionAPIPath + "/v1/deputies/%d/notes", deputyId), &body)

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

	if resp.StatusCode != http.StatusCreated {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{Errors: v.ValidationErrors}
		}

		return newStatusError(resp)
	}

	return nil
}

func getNoteType(deputyType string) string {
	if deputyType == "PRO" {
		return "PRO_DEPUTY_NOTE_CREATED"
	} else {
		return "PA_DEPUTY_NOTE_CREATED"
	}
}
