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

	body, _ = CreateAddDocumentFormBody(body, writer, "GENERAL", "OUTGOING", "01/01/2024", "", false)

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
	body, _ = CreateAddDocumentFormBody(body, writer, "GENERAL", "OUTGOING", "01/01/2024", "", false)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	returnedError := renderTemplateForAddDocument(client, template)(app, w, r)

	assert.Equal(client.AddDocumentErr, returnedError)
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
			body, _ = CreateAddDocumentFormBody(body, writer, "GENERAL", "OUTGOING", "01/01/2024", "", false)
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/123", body)
			r.Header.Add("Content-Type", writer.FormDataContentType())
			returnedError := renderTemplateForAddDocument(client, template)(app, w, r)

			refDataReturnedError := renderTemplateForAddDocument(client, template)(app, w, r)
			assert.Equal(t, returnedError, refDataReturnedError)
		})
	}
}

func TestAddDocumentHandlesFileUploadError(t *testing.T) {
	assert := assert.New(t)

	expectedError := sirius.ValidationErrors{
		"document-upload": {
			"": "Error uploading the file",
		},
	}

	client := &mockAddDocumentClient{}
	template := &mockTemplates{}
	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	body, _ = CreateAddDocumentFormBody(body, writer, "GENERAL", "OUTGOING", "01/01/2024", "", true)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	returnedError := renderTemplateForAddDocument(client, template)(app, w, r)

	assert.Equal(AddDocumentVars{
		AppVars: AppVars{
			DeputyDetails: sirius.DeputyDetails{ID: 123},
			Errors:        expectedError,
			PageName:      "Add a document",
			Path:          "/path",
		},
		DocumentDirectionRefData: []model.RefData{},
		DocumentTypes:            []model.RefData{},
		DocumentType:             "GENERAL",
		Direction:                "OUTGOING",
		Date:                     "01/01/2024",
		Notes:                    "",
	}, template.lastVars)

	assert.Nil(returnedError)
}

func CreateAddDocumentFormBody(body *bytes.Buffer, writer *multipart.Writer, documentType, direction, date, notes string, documentUploadError bool) (*bytes.Buffer, error) {
	if documentUploadError {
		_, _ = writer.CreateFormFile("document-upload", "")
	} else {
		dataPart, _ := writer.CreateFormFile("document-upload", "data.txt")
		_, _ = io.Copy(dataPart, strings.NewReader("blarg"))
	}

	if documentType != "" {
		typeWriter, err := writer.CreateFormField("documentType")
		if err != nil {
			return nil, err
		}
		_, err = typeWriter.Write([]byte(documentType))
		if err != nil {
			return nil, err
		}
	}

	if direction != "" {
		directionWriter, err := writer.CreateFormField("documentDirection")
		if err != nil {
			return nil, err
		}
		_, err = directionWriter.Write([]byte(direction))
		if err != nil {
			return nil, err
		}
	}

	if date != "" {
		dateWriter, err := writer.CreateFormField("documentDate")
		if err != nil {
			return nil, err
		}
		_, err = dateWriter.Write([]byte(date))
		if err != nil {
			return nil, err
		}
	}

	notesWriter, err := writer.CreateFormField("notes")
	if err != nil {
		return nil, err
	}
	_, err = notesWriter.Write([]byte(notes))
	if err != nil {
		return nil, err
	}

	writer.Close()
	return body, nil
}
