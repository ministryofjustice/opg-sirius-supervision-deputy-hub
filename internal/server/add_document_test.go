package server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddDocumentClient struct {
	count                       int
	lastCtx                     sirius.Context
	AddDocumentErr              error
	GetDocumentTypesRefData     error
	GetDocumentDirectionRefData error
}

func (m *mockAddDocumentClient) AddDocument(ctx sirius.Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.AddDocumentErr
}

func (m *mockAddDocumentClient) GetRefData(ctx sirius.Context, refDataUrlType string) ([]model.RefData, error) {
	m.count += 1
	m.lastCtx = ctx

	return []model.RefData{}, m.GetDocumentTypesRefData
}

var addDocumentVars = AppVars{
	DeputyDetails: sirius.DeputyDetails{
		ID:              123,
		DeputyFirstName: "Test",
		DeputySurname:   "Dep",
	},
}

func TestGetAddDocument(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddDocumentClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/123", strings.NewReader(""))

	handler := renderTemplateForAddDocument(client, template)
	err := handler(addDocumentVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostAddDocument(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddDocumentClient{}
	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	body, _ = CreateAddDocumentFormBody(body, writer, true, true, true, false, false)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	var res error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		res = renderTemplateForAddDocument(client, nil)(app, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(res, Redirect(fmt.Sprintf("/123/documents?success=addDocument&filename=%s", "data.txt")))
}

func TestPostAddDocumentReturnsErrorsFromSirius(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddDocumentClient{
		AddDocumentErr: sirius.StatusError{Code: 500},
	}
	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
	}

	template := &mockTemplates{}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	body, _ = CreateAddDocumentFormBody(body, writer, true, true, true, false, false)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	returnedError := renderTemplateForAddDocument(client, template)(app, w, r)

	assert.Equal(client.AddDocumentErr, returnedError)
}

func TestAddDocumentHandlesValidationErrorsInternally(t *testing.T) {
	tests := []struct {
		ExpectedError          sirius.ValidationErrors
		hasType                bool
		hasDirection           bool
		hasDate                bool
		hasNotesError          bool
		hasDocumentUploadError bool
	}{
		{
			ExpectedError: sirius.ValidationErrors{
				"type": {
					"": "Select a type",
				},
			},
			hasType:                false,
			hasDirection:           true,
			hasDate:                true,
			hasNotesError:          false,
			hasDocumentUploadError: false,
		},
		{
			ExpectedError: sirius.ValidationErrors{
				"direction": {
					"": "Select a direction",
				},
			},
			hasType:                true,
			hasDirection:           false,
			hasDate:                true,
			hasNotesError:          false,
			hasDocumentUploadError: false,
		},
		{
			ExpectedError: sirius.ValidationErrors{
				"date": {
					"": "Select a date",
				},
			},
			hasType:                true,
			hasDirection:           true,
			hasDate:                false,
			hasNotesError:          false,
			hasDocumentUploadError: false,
		},
		{
			ExpectedError: sirius.ValidationErrors{
				"notes": {
					"stringLengthTooLong": "The note must be 1000 characters or fewer",
				},
			},
			hasType:                true,
			hasDirection:           true,
			hasDate:                true,
			hasNotesError:          true,
			hasDocumentUploadError: false,
		},
		{
			ExpectedError: sirius.ValidationErrors{
				"document-upload": {
					"": "Error uploading the file",
				},
			},
			hasType:                true,
			hasDirection:           true,
			hasDate:                true,
			hasNotesError:          false,
			hasDocumentUploadError: true,
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {
			assert := assert.New(t)

			client := &mockAddDocumentClient{}
			template := &mockTemplates{}
			app := AppVars{
				Path:          "/path",
				DeputyDetails: sirius.DeputyDetails{ID: 123},
			}
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			body, _ = CreateAddDocumentFormBody(body, writer, tc.hasType, tc.hasDirection, tc.hasDate, tc.hasNotesError, tc.hasDocumentUploadError)
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/123", body)
			r.Header.Add("Content-Type", writer.FormDataContentType())
			returnedError := renderTemplateForAddDocument(client, template)(app, w, r)

			assert.Equal(AddDocumentVars{
				AppVars: AppVars{
					DeputyDetails: sirius.DeputyDetails{ID: 123},
					Errors:        tc.ExpectedError,
					PageName:      "Add a document",
					Path:          "/path",
				},
				DocumentDirectionRefData: []model.RefData{},
				DocumentTypes:            []model.RefData{},
			}, template.lastVars)

			assert.Nil(returnedError)
		})
	}
}

func TestAddDocumentHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockAddDocumentClient
	}{
		{
			Client: &mockAddDocumentClient{
				GetDocumentDirectionRefData: returnedError,
			},
		},
		{
			Client: &mockAddDocumentClient{
				GetDocumentTypesRefData: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}
			app := AppVars{
				Path:          "/path",
				DeputyDetails: sirius.DeputyDetails{ID: 123},
			}
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			body, _ = CreateAddDocumentFormBody(body, writer, true, true, true, false, false)
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/123", body)
			r.Header.Add("Content-Type", writer.FormDataContentType())
			returnedError := renderTemplateForAddDocument(client, template)(app, w, r)

			refDataReturnedError := renderTemplateForAddDocument(client, template)(app, w, r)
			assert.Equal(t, returnedError, refDataReturnedError)
		})
	}
}

func CreateAddDocumentFormBody(body *bytes.Buffer, writer *multipart.Writer, hasType, hasDirection, hasDate, hasNotesError, documentUploadError bool) (*bytes.Buffer, error) {

	if documentUploadError {
		_, _ = writer.CreateFormFile("document-upload", "")
	} else {
		dataPart, _ := writer.CreateFormFile("document-upload", "data.txt")
		_, _ = io.Copy(dataPart, strings.NewReader("blarg"))
	}

	if hasType {
		typeWriter, err := writer.CreateFormField("type")
		if err != nil {
			return body, err
		}
		_, err = typeWriter.Write([]byte("ABC"))
		if err != nil {
			return nil, err
		}
	}

	if hasDirection {
		direction, err := writer.CreateFormField("direction")
		if err != nil {
			return body, err
		}
		_, err = direction.Write([]byte("INCOMING"))
		if err != nil {
			return nil, err
		}
	}

	if hasDate {
		date, err := writer.CreateFormField("date")
		if err != nil {
			return body, err
		}
		_, err = date.Write([]byte("2020-01-01"))
		if err != nil {
			return nil, err
		}
	}

	notes, err := writer.CreateFormField("notes")
	if err != nil {
		return body, err
	}

	count := 1
	if hasNotesError {
		count = 101
	}

	content := strings.Repeat("1234567890", count)
	_, err = notes.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	writer.Close()
	return body, nil
}
