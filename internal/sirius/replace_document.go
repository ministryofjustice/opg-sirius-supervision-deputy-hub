package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"mime/multipart"
	"net/http"
	"time"
)

type T struct {
	Class   string `json:"class"`
	Payload struct {
		IsPersonAndCaseEvent bool      `json:"isPersonAndCaseEvent"`
		IsPersonEvent        bool      `json:"isPersonEvent"`
		IsCaseEvent          bool      `json:"isCaseEvent"`
		DocumentId           string    `json:"documentId"`
		Filename             string    `json:"filename"`
		Description          string    `json:"description"`
		CreatedBy            string    `json:"createdBy"`
		CreatedDate          time.Time `json:"createdDate"`
		Direction            string    `json:"direction"`
		ReceivedDate         time.Time `json:"receivedDate"`
		Reason               string    `json:"reason"`
		Type                 string    `json:"type"`
		PersonType           string    `json:"personType"`
		PersonId             string    `json:"personId"`
		PersonUid            string    `json:"personUid"`
		PersonName           string    `json:"personName"`
		Changes              []struct {
			FieldName string      `json:"fieldName"`
			OldValue  string      `json:"oldValue,omitempty"`
			NewValue  interface{} `json:"newValue"`
			Type      string      `json:"type"`
		} `json:"changes"`
	} `json:"payload"`
}

func (c *Client) ReplaceDocument(ctx Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId, documentId int) error {
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
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/deputies/%d/documents/%d/replace-file", deputyId, documentId), &body)

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
