package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"mime/multipart"
	"net/http"
)

type CreateDocumentRequest struct {
	Date        string        `json:"documentDate"`
	Description string        `json:"description"`
	Direction   model.RefData `json:"documentDirection"`
	Name        string        `json:"name"`
	Type        model.RefData `json:"documentType"`
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

	source, err := model.EncodeFileToBase64(file)
	if err != nil {
		return err
	}

	requestBody := CreateDocumentRequest{
		File: EncodedFile{
			Name:   filename,
			Source: source,
			Type:   "application/json",
		},
		FileName:   filename,
		FileSource: source,
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

	err = json.NewEncoder(&body).Encode(requestBody)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/documents", deputyId), &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)

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
