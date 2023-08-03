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

type mockAddAssuranceVisitInformation struct {
	count                int
	lastCtx              sirius.Context
	AddAssuranceVisitErr error
	err                  error
}

func (m *mockAddAssuranceVisitInformation) AddAssuranceVisit(ctx sirius.Context, assuranceType string, requestedDate string, userId, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.AddAssuranceVisitErr
}

func TestGetAddAssuranceVisit(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddAssuranceVisitInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddAssuranceVisit(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostAssuranceVisit(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddAssuranceVisitInformation{}
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
		returnedError = renderTemplateForAddAssuranceVisit(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(Redirect("/123/assurance-visits?success=addAssuranceVisit"), returnedError)
}

func TestAssuranceVisitHandlesValidationErrorsGeneratedWithinFile(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddAssuranceVisitInformation{}

	form := url.Values{}
	form.Add("assurance-type", "")
	form.Add("requested-date", "")

	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurance-visits", strings.NewReader(form.Encode()))

	returnedError := renderTemplateForAddAssuranceVisit(client, template)(AppVars{}, w, r)

	expectedErrors := sirius.ValidationErrors{
		"assurance-type": {
			"": "Select an assurance type",
		},
		"requested-date": {
			"": "Enter a requested date",
		},
	}

	assert.Equal(AddAssuranceVisitVars{
		AppVars{
			Path:   "/123/assurance-visits",
			Errors: expectedErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestAssuranceVisitHandlesValidationErrorsReturnedFromSiriusCall(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddAssuranceVisitInformation{}

	validationErrors := sirius.ValidationErrors{
		"assurance-type": {
			"": "Select an assurance type",
		},
		"requested-date": {
			"": "Enter a requested date",
		},
	}

	client.AddAssuranceVisitErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurance-visits", strings.NewReader(""))

	returnedError := renderTemplateForAddAssuranceVisit(client, template)(AppVars{}, w, r)

	assert.Equal(AddAssuranceVisitVars{
		AppVars{
			Path:   "/123/assurance-visits",
			Errors: validationErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)

}
