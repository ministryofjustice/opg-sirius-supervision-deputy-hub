package sirius

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"io"
	"mime/multipart"
	"net/http"
)

type CreateDocumentRequest struct {
	Date        string        `json:"date"`
	Description string        `json:"description"`
	Direction   model.RefData `json:"direction"`
	Name        string        `json:"name"`
	Type        model.RefData `json:"type"`
	PersonId    int           `json:"personId"`
	FileName    string        `json:"fileName"`
	FileSource  string        `json:"fileSource"`
	File        EncodedFile   `json:"file"`
}

type EncodedFile struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

func (c *Client) AddDocument(ctx Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId int) error {
	var body bytes.Buffer

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		fmt.Println(err)
	}

	source := base64.StdEncoding.EncodeToString(buf.Bytes())

	requestBody := CreateDocumentRequest{
		File: EncodedFile{
			Name:   filename,
			Source: source,
			Type:   "application/json",
		},
		FileName:   filename,
		FileSource: source,
		PersonId:   deputyId,
		Date:       date,
		Type: model.RefData{
			Handle: documentType,
		},
		Direction: model.RefData{
			Handle: direction,
		},
		Name:        "Document Uploaded",
		Description: notes,
	}

	file.Close()

	err := json.NewEncoder(&body).Encode(requestBody)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/api/public/v1/documents/deputies", &body)

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
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
