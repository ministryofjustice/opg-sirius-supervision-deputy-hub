package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddDocumentClient struct {
	count          int
	lastCtx        sirius.Context
	refData        []model.RefData
	addDocumentErr error
	getRefDataErr  error
	successMessage string
}

func (m *mockAddDocumentClient) AddDocument(ctx sirius.Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.addDocumentErr
}

func (m *mockAddDocumentClient) GetRefData(ctx sirius.Context, refDataUrlType string) ([]model.RefData, error) {
	m.count += 1
	m.lastCtx = ctx

	refData := []model.RefData{
		{
			Handle: "HANDLE",
			Label:  "label",
		},
	}

	return refData, m.getRefDataErr
}

var addDocumentVars = AddDocumentVars{
	AppVars: AppVars{},
}

func TestAddDocument(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddDocumentClient{}

	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
	}

	//body := new(bytes.Buffer)
	//writer := multipart.NewWriter(body)
	//// create a new form-data header name data and filename data.txt
	//dataPart, _ := writer.CreateFormFile("data", "data.txt")
	//_, _ = io.CopyBuffer(dataPart, []byte("test"))

	//var mockForm *multipart.Form
	//mockForm

	form := url.Values{
		"type":            {"ABC"},
		"direction":       {"INCOMING"},
		"date":            {"2020-01-01"},
		"notes":           {"Notes on this file"},
		"document-upload": {"some content"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(form.Encode()))

	//var v any
	//json.NewDecoder(body).Decode(&v)
	//fmt.Print(v)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var res error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		res = renderTemplateForAddDocument(client, nil)(app, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(res, Redirect("123/documents?success=addDocument&filename=%s"))
}

//
//func TestPostAddDocument(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockAddDocumentClient{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForAddFirm(client, nil)(addFirmAppVars, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//	assert.Equal(returnedError, Redirect("/123?success=newFirm"))
//}

//func TestAddFirmValidationErrors(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockFirmInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"stringLengthTooLong": "The firm name must be 255 characters or fewer",
//		},
//	}
//
//	client.AddFirmDetailsErr = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
//	returnedError := renderTemplateForAddFirm(client, template)(AppVars{}, w, r)
//
//	assert.Equal(addFirmVars{
//		AppVars: AppVars{
//			Errors:   validationErrors,
//			PageName: "Create new firm",
//		},
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}
//
//func TestErrorAddFirmMessageWhenIsEmpty(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockFirmInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"isEmpty": "The firm name is required and can't be empty",
//		},
//	}
//
//	client.AddFirmDetailsErr = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForAddFirm(client, template)(addFirmAppVars, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//
//	expectedValidationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"isEmpty": "The firm name is required and can't be empty",
//		},
//	}
//
//	assert.Equal(addFirmVars{
//		AppVars: AppVars{
//			DeputyDetails: testDeputy,
//			Errors:        expectedValidationErrors,
//			PageName:      "Create new firm",
//		},
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}
//
//func TestAddFirmHandlesErrorsInOtherClientFiles(t *testing.T) {
//	returnedError := sirius.StatusError{Code: 500}
//	tests := []struct {
//		Client *mockFirmInformation
//	}{
//		{
//			Client: &mockFirmInformation{
//				AddFirmDetailsErr: returnedError,
//			},
//		},
//		{
//			Client: &mockFirmInformation{
//				AssignDeputyToFirmErr: returnedError,
//			},
//		},
//	}
//	for k, tc := range tests {
//		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {
//
//			client := tc.Client
//			template := &mockTemplates{}
//			w := httptest.NewRecorder()
//			r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
//
//			addFirmReturnedError := renderTemplateForAddFirm(client, template)(AppVars{}, w, r)
//			assert.Equal(t, returnedError, addFirmReturnedError)
//		})
//	}
//}
