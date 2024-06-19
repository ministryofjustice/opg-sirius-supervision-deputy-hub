package sirius

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type CreateDocument struct {
	Type          string `json:"assuranceType"`
	RequestedDate string `json:"requestedDate"`
	RequestedBy   int    `json:"requestedBy"`
}

type TestFile struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

type AddDocumentRequest struct {
	Type          string   `json:"type"`
	CaseRecNumber string   `json:"caseRecNumber"`
	ParentUuid    string   `json:"parentUuid"`
	Metadata      string   `json:"metadata"`
	File          TestFile `json:"file"`
}

type CreateNote struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	Direction   struct {
		Handle string `json:"handle"`
		Label  string `json:"label"`
	} `json:"direction"`
	Name string `json:"name"`
	Type struct {
		Handle string `json:"handle"`
		Label  string `json:"label"`
	} `json:"type"`
	PersonId   int    `json:"personId"`
	FileName   string `json:"fileName"`
	FileSource string `json:"fileSource"`
	File       struct {
		Name   string `json:"name"`
		Source string `json:"source"`
		Type   string `json:"type"`
	} `json:"file"`
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

	fmt.Print("documentType")
	fmt.Println(documentType)

	source := base64.StdEncoding.EncodeToString(buf.Bytes())

	requestBody := CreateNote{
		File: EncodedFile{
			Name:   filename,
			Source: source,
			Type:   "application/json",
		},
		FileName:   filename,
		FileSource: source,
		PersonId:   deputyId,
		Date:       "01/01/2020",
		Type: struct {
			Handle string `json:"handle"`
			Label  string `json:"label"`
		}{
			Handle: documentType,
		},
		Direction: struct {
			Handle string `json:"handle"`
			Label  string `json:"label"`
		}{
			Handle: direction,
		},
	}

	//date is received date time (I think)

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

	//io.Copy(os.Stdout, resp.Body)

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
