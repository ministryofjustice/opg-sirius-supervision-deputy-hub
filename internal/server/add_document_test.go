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

func (m *mockAddDocumentClient) GetDocumentTypes(ctx sirius.Context) ([]model.RefData, error) {
	m.count += 1
	m.lastCtx = ctx

	return []model.RefData{}, m.GetDocumentTypesRefData
}

func (m *mockAddDocumentClient) GetDocumentDirections(ctx sirius.Context) ([]model.RefData, error) {
	m.count += 1
	m.lastCtx = ctx

	return []model.RefData{}, m.GetDocumentDirectionRefData
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

	body, _ = CreateDocumentFormBody(body, writer, false)

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
	body, _ = CreateDocumentFormBody(body, writer, false)

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
			body, _ = CreateDocumentFormBody(body, writer, false)
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
			"": "Select a file to attach",
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
	body, _ = CreateDocumentFormBody(body, writer, true)

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
		DocumentTypes:            []model.RefData(nil),
		DocumentType:             "GENERAL",
		Direction:                "OUTGOING",
		Date:                     "01/01/2024",
		Notes:                    "",
	}, template.lastVars)

	assert.Nil(returnedError)
}

func CreateDocumentFormBody(body *bytes.Buffer, writer *multipart.Writer, documentUploadError bool) (*bytes.Buffer, error) {
	if documentUploadError {
		_, _ = writer.CreateFormFile("document-upload", "")
	} else {
		dataPart, _ := writer.CreateFormFile("document-upload", "data.txt")
		_, _ = io.Copy(dataPart, strings.NewReader("blarg"))
	}

	typeWriter, err := writer.CreateFormField("documentType")
	if err != nil {
		return nil, err
	}
	_, err = typeWriter.Write([]byte("GENERAL"))
	if err != nil {
		return nil, err
	}

	directionWriter, err := writer.CreateFormField("documentDirection")
	if err != nil {
		return nil, err
	}
	_, err = directionWriter.Write([]byte("OUTGOING"))
	if err != nil {
		return nil, err
	}

	dateWriter, err := writer.CreateFormField("documentDate")
	if err != nil {
		return nil, err
	}
	_, err = dateWriter.Write([]byte("01/01/2024"))
	if err != nil {
		return nil, err
	}

	notesWriter, err := writer.CreateFormField("notes")
	if err != nil {
		return nil, err
	}
	_, err = notesWriter.Write([]byte(""))
	if err != nil {
		return nil, err
	}

	writer.Close()
	return body, nil
}

func Test_filterDocTypeByDeputyType(t *testing.T) {
	documentTypes := []model.RefData{{Handle: "ASSURANCE_VISIT", Label: "Assurance visit", Deprecated: false}, {Handle: "CATCH_UP_CALL", Label: "Catch-up call", Deprecated: false}, {Handle: "COMPLAINTS", Label: "Complaints", Deprecated: false}, {Handle: "CORRESPONDENCE", Label: "Correspondence", Deprecated: false}, {Handle: "GENERAL", Label: "General", Deprecated: false}, {Handle: "INDEMNITY_INSURANCE", Label: "Indemnity insurance", Deprecated: false}, {Handle: "NON_COMPLIANCE", Label: "Non-compliance", Deprecated: false}}
	tests := []struct {
		name          string
		DocumentTypes []model.RefData
		DeputyType    string
		want          []model.RefData
	}{
		{
			name:          "Document list types returns without catch-up call for a PRO Deputy",
			DocumentTypes: documentTypes,
			DeputyType:    "PRO",
			want:          []model.RefData{{Handle: "ASSURANCE_VISIT", Label: "Assurance visit", Deprecated: false}, {Handle: "COMPLAINTS", Label: "Complaints", Deprecated: false}, {Handle: "CORRESPONDENCE", Label: "Correspondence", Deprecated: false}, {Handle: "GENERAL", Label: "General", Deprecated: false}, {Handle: "INDEMNITY_INSURANCE", Label: "Indemnity insurance", Deprecated: false}, {Handle: "NON_COMPLIANCE", Label: "Non-compliance", Deprecated: false}},
		},
		{
			name:          "Document list types returns without Indemnity insurance for a PA Deputy",
			DocumentTypes: documentTypes,
			DeputyType:    "PA",
			want:          []model.RefData{{Handle: "ASSURANCE_VISIT", Label: "Assurance visit", Deprecated: false}, {Handle: "CATCH_UP_CALL", Label: "Catch-up call", Deprecated: false}, {Handle: "COMPLAINTS", Label: "Complaints", Deprecated: false}, {Handle: "CORRESPONDENCE", Label: "Correspondence", Deprecated: false}, {Handle: "GENERAL", Label: "General", Deprecated: false}, {Handle: "NON_COMPLIANCE", Label: "Non-compliance", Deprecated: false}},
		},
		{
			name:          "Document list types returns all with a deputy type",
			DocumentTypes: documentTypes,
			DeputyType:    "",
			want:          []model.RefData(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, filterDocTypeByDeputyType(tt.DocumentTypes, tt.DeputyType), "filterDocTypeByDeputyType(%v, %v)", tt.DocumentTypes, tt.DeputyType)
		})
	}
}
