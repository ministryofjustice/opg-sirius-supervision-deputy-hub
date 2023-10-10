package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddAssuranceClient struct {
	count           int
	lastCtx         sirius.Context
	AddAssuranceErr error
}

func (m *mockAddAssuranceClient) AddAssurance(ctx sirius.Context, assuranceType string, requestedDate string, userId, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.AddAssuranceErr
}

func TestGetAddAssurance(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddAssuranceClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddAssurance(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostAssurance(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddAssuranceClient{}
	template := &mockTemplates{}

	form := url.Values{}
	form.Add("assurance-type", "ABC")
	form.Add("requested-date", "2200/10/20")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurance-visits", strings.NewReader(form.Encode()))
	r.PostForm = form

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/assurance-visits", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddAssurance(client, template)(AppVars{DeputyDetails: testDeputy}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(Redirect("/123/assurance-visits?success=addAssuranceVisit"), returnedError)
}

func TestAddAssuranceHandlesValidationErrorsGeneratedWithinFile(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddAssuranceClient{}

	form := url.Values{}
	form.Add("assurance-type", "")
	form.Add("requested-date", "")

	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurance-visits", strings.NewReader(form.Encode()))

	returnedError := renderTemplateForAddAssurance(client, template)(AppVars{}, w, r)

	expectedErrors := sirius.ValidationErrors{
		"assurance-type": {
			"": "Select an assurance type",
		},
		"requested-date": {
			"": "Enter a requested date",
		},
	}

	assert.Equal(AddAssuranceVars{
		AppVars{
			Errors: expectedErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestAddAssuranceHandlesValidationErrorsReturnedFromSiriusCall(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddAssuranceClient{}

	validationErrors := sirius.ValidationErrors{
		"assurance-type": {
			"": "Select an assurance type",
		},
		"requested-date": {
			"": "Enter a requested date",
		},
	}

	client.AddAssuranceErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurance-visits", strings.NewReader(""))

	returnedError := renderTemplateForAddAssurance(client, template)(AppVars{}, w, r)

	assert.Equal(AddAssuranceVars{
		AppVars{
			Errors: validationErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)

}
