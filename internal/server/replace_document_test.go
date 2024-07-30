package server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

type mockReplaceDocumentClient struct {
	count                       int
	lastCtx                     sirius.Context
	DocumentToReplace           model.Document
	ReplaceDocumentErr          error
	GetDocumentByIdErr          error
	GetDocumentTypesRefData     error
	GetDocumentDirectionRefData error
}

func (m *mockReplaceDocumentClient) ReplaceDocument(ctx sirius.Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId, documentId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.ReplaceDocumentErr
}

func (m *mockReplaceDocumentClient) GetDocumentDirections(ctx sirius.Context) ([]model.RefData, error) {
	m.count += 1
	m.lastCtx = ctx

	return []model.RefData{}, m.GetDocumentDirectionRefData
}

func (m *mockReplaceDocumentClient) GetDocumentTypes(ctx sirius.Context) ([]model.RefData, error) {
	m.count += 1
	m.lastCtx = ctx

	return []model.RefData{}, m.GetDocumentTypesRefData
}

func (m *mockReplaceDocumentClient) GetDocumentById(ctx sirius.Context, deputyId, documentId int) (model.Document, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.DocumentToReplace, m.GetDocumentByIdErr
}

var replaceDocumentVars = AppVars{
	DeputyDetails: sirius.DeputyDetails{
		ID:              123,
		DeputyFirstName: "Test",
		DeputySurname:   "Dep",
	},
	Path: "/deputies/123/replace",
}

func TestGetReplaceDocument(t *testing.T) {
	assert := assert.New(t)
	client := &mockReplaceDocumentClient{
		DocumentToReplace: model.Document{
			Id:                  5,
			Type:                "GENERAL",
			FriendlyDescription: "bad-doc.png",
			Filename:            "1234_bad-doc.png",
			ReceivedDateTime:    "04/07/2024 01:00:00",
			Direction:           "INCOMING",
		},
	}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "deputies/123/documents/5/replace", strings.NewReader(""))

	handler := renderTemplateForReplaceDocument(client, template)
	err := handler(replaceDocumentVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(ReplaceDocumentVars{
		AppVars: AppVars{
			DeputyDetails: sirius.DeputyDetails{
				ID:              123,
				DeputyFirstName: "Test",
				DeputySurname:   "Dep",
			},
			PageName: "Replace a document",
			Path:     "/deputies/123/replace",
		},
		OriginalDocument: model.Document{
			Id:                  5,
			Type:                "GENERAL",
			FriendlyDescription: "bad-doc.png",
			CreatedDate:         "",
			Direction:           "INCOMING",
			Filename:            "1234_bad-doc.png",
			CreatedBy:           model.User{},
			ReceivedDateTime:    "04/07/2024 01:00:00",
			ReformattedTime:     "04/07/2024",
			Note:                model.DocumentNote{},
		},
		DocumentDirectionRefData: []model.RefData{},
		DocumentTypes:            []model.RefData{},
		Date:                     time.Now().Format("2006-01-02"),
		Notes:                    "",
	}, template.lastVars)
}

func TestPostReplaceDocument(t *testing.T) {
	assert := assert.New(t)

	client := &mockReplaceDocumentClient{
		DocumentToReplace: model.Document{
			Id:                  5,
			Type:                "GENERAL",
			FriendlyDescription: "bad-doc.png",
			Filename:            "1234_bad-doc.png",
			ReceivedDateTime:    "04/07/2024 01:00:00",
			Direction:           "INCOMING",
		},
	}
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
		res = renderTemplateForReplaceDocument(client, nil)(app, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(res, Redirect(fmt.Sprintf("/123/documents?success=replaceDocument&previousFilename=%s&filename=%s", "bad-doc.png", "data.txt")))
}

func TestPostReplaceDocumentReturnsErrorsFromSirius(t *testing.T) {
	assert := assert.New(t)
	client := &mockReplaceDocumentClient{
		DocumentToReplace: model.Document{
			Id:                  5,
			Type:                "GENERAL",
			FriendlyDescription: "bad-doc.png",
			Filename:            "1234_bad-doc.png",
			ReceivedDateTime:    "04/07/2024 01:00:00",
			Direction:           "INCOMING",
		},
		ReplaceDocumentErr: sirius.StatusError{Code: 500},
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
	returnedError := renderTemplateForReplaceDocument(client, template)(app, w, r)

	assert.Equal(client.ReplaceDocumentErr, returnedError)
}

func TestReplaceDocumentHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockReplaceDocumentClient
	}{
		{
			Client: &mockReplaceDocumentClient{
				GetDocumentDirectionRefData: returnedError,
				DocumentToReplace: model.Document{
					Id:                  5,
					Type:                "GENERAL",
					FriendlyDescription: "bad-doc.png",
					Filename:            "1234_bad-doc.png",
					ReceivedDateTime:    "04/07/2024 01:00:00",
					Direction:           "INCOMING",
				},
			},
		},
		{
			Client: &mockReplaceDocumentClient{
				GetDocumentTypesRefData: returnedError,
				DocumentToReplace: model.Document{
					Id:                  5,
					Type:                "GENERAL",
					FriendlyDescription: "bad-doc.png",
					Filename:            "1234_bad-doc.png",
					ReceivedDateTime:    "04/07/2024 01:00:00",
					Direction:           "INCOMING",
				},
			},
		},
		{
			Client: &mockReplaceDocumentClient{
				GetDocumentByIdErr: returnedError,
				DocumentToReplace: model.Document{
					Id:                  5,
					Type:                "GENERAL",
					FriendlyDescription: "bad-doc.png",
					Filename:            "1234_bad-doc.png",
					ReceivedDateTime:    "04/07/2024 01:00:00",
					Direction:           "INCOMING",
				},
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
			returnedError := renderTemplateForReplaceDocument(client, template)(app, w, r)

			refDataReturnedError := renderTemplateForReplaceDocument(client, template)(app, w, r)
			assert.Equal(t, returnedError, refDataReturnedError)
		})
	}
}

func TestReplaceDocumentHandlesFileUploadError(t *testing.T) {
	assert := assert.New(t)

	expectedError := sirius.ValidationErrors{
		"document-upload": {
			"": "Select a file to attach",
		},
	}

	client := &mockReplaceDocumentClient{
		DocumentToReplace: model.Document{
			Id:                  5,
			Type:                "GENERAL",
			FriendlyDescription: "bad-doc.png",
			Filename:            "1234_bad-doc.png",
			ReceivedDateTime:    "04/07/2024 01:00:00",
			Direction:           "INCOMING",
		},
	}
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
	returnedError := renderTemplateForReplaceDocument(client, template)(app, w, r)

	assert.Equal(ReplaceDocumentVars{
		AppVars: AppVars{
			DeputyDetails: sirius.DeputyDetails{ID: 123},
			Errors:        expectedError,
			PageName:      "Replace a document",
			Path:          "/path",
		},
		OriginalDocument: model.Document{
			Id:                  5,
			Type:                "GENERAL",
			FriendlyDescription: "bad-doc.png",
			CreatedDate:         "",
			Direction:           "INCOMING",
			Filename:            "1234_bad-doc.png",
			CreatedBy:           model.User{},
			ReceivedDateTime:    "04/07/2024 01:00:00",
			ReformattedTime:     "04/07/2024",
			Note:                model.DocumentNote{},
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
