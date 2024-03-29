package server

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockFirmInformation struct {
	count                 int
	lastCtx               sirius.Context
	AddFirmDetailsErr     error
	AssignDeputyToFirmErr error
	addFirm               int
}

func (m *mockFirmInformation) AddFirmDetails(ctx sirius.Context, deputyId sirius.FirmDetails) (int, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.addFirm, m.AddFirmDetailsErr
}

func (m *mockFirmInformation) AssignDeputyToFirm(ctx sirius.Context, deputyId int, firmId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.AssignDeputyToFirmErr
}

var addFirmAppVars = AppVars{
	DeputyDetails: testDeputy,
}

func TestGetFirm(t *testing.T) {
	assert := assert.New(t)

	client := &mockFirmInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddFirm(client, template)
	err := handler(addFirmAppVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostAddFirm(t *testing.T) {
	assert := assert.New(t)
	client := &mockFirmInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddFirm(client, nil)(addFirmAppVars, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/123?success=newFirm"))
}

func TestAddFirmValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockFirmInformation{}

	validationErrors := sirius.ValidationErrors{
		"firmName": {
			"stringLengthTooLong": "The firm name must be 255 characters or fewer",
		},
	}

	client.AddFirmDetailsErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
	returnedError := renderTemplateForAddFirm(client, template)(AppVars{}, w, r)

	assert.Equal(addFirmVars{
		AppVars: AppVars{
			Errors:   validationErrors,
			PageName: "Create new firm",
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestErrorAddFirmMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockFirmInformation{}

	validationErrors := sirius.ValidationErrors{
		"firmName": {
			"isEmpty": "The firm name is required and can't be empty",
		},
	}

	client.AddFirmDetailsErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddFirm(client, template)(addFirmAppVars, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"firmName": {
			"isEmpty": "The firm name is required and can't be empty",
		},
	}

	assert.Equal(addFirmVars{
		AppVars: AppVars{
			DeputyDetails: testDeputy,
			Errors:        expectedValidationErrors,
			PageName:      "Create new firm",
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestAddFirmHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockFirmInformation
	}{
		{
			Client: &mockFirmInformation{
				AddFirmDetailsErr: returnedError,
			},
		},
		{
			Client: &mockFirmInformation{
				AssignDeputyToFirmErr: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

			addFirmReturnedError := renderTemplateForAddFirm(client, template)(AppVars{}, w, r)
			assert.Equal(t, returnedError, addFirmReturnedError)
		})
	}
}
