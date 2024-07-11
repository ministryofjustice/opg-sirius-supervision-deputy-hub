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

func EncodeFileToBase64(file multipart.File) (string, error) {
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	file.Close()

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (c *Client) AddDocument(ctx Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId int) error {
	var body bytes.Buffer

	source, err := EncodeFileToBase64(file)
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

	err = json.NewEncoder(&body).Encode(requestBody)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/deputies/%d/documents", deputyId), &body)

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
