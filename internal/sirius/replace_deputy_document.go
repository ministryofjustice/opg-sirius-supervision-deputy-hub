package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"mime/multipart"
	"net/http"
)

func (c *Client) ReplaceDocument(ctx Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId, documentId int) error {
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
		FileName: filename,
		Date:     date,
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
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/deputies/%d/documents/%d", deputyId, documentId), &body)

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
