package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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
	GetUserDetailsErr    error
	userDetails          sirius.UserDetails
}

func (m *mockAddAssuranceVisitInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetails, m.GetUserDetailsErr
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
	err := handler(sirius.DeputyDetails{}, w, r)

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
		returnedError = renderTemplateForAddAssuranceVisit(client, template)(sirius.DeputyDetails{}, w, r)
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

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/assurance-visits", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddAssuranceVisit(client, template)(sirius.DeputyDetails{}, w, r)
	})
	testHandler.ServeHTTP(w, r)

	expectedErrors := sirius.ValidationErrors{
		"assurance-type": {
			"": "Select an assurance type",
		},
		"requested-date": {
			"": "Enter a requested date",
		},
	}

	assert.Equal(AddAssuranceVisitVars{
		Path:   "/123/assurance-visits",
		Errors: expectedErrors,
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

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/assurance-visits", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddAssuranceVisit(client, template)(sirius.DeputyDetails{}, w, r)
	})
	testHandler.ServeHTTP(w, r)

	assert.Equal(AddAssuranceVisitVars{
		Path:   "/123/assurance-visits",
		Errors: validationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)

}

func TestAddAssuranceVisitInformationHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockAddAssuranceVisitInformation
	}{
		{
			Client: &mockAddAssuranceVisitInformation{
				GetUserDetailsErr: returnedError,
			},
		},
		{
			Client: &mockAddAssuranceVisitInformation{
				AddAssuranceVisitErr: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}

			w := httptest.NewRecorder()
			form := url.Values{}
			form.Add("assurance-type", "ABC")
			form.Add("requested-date", "2200/10/20")

			r, _ := http.NewRequest("POST", "/123/assurance-visits", strings.NewReader(form.Encode()))
			r.PostForm = form

			addAssuranceVisitReturnedError := renderTemplateForAddAssuranceVisit(client, template)(sirius.DeputyDetails{}, w, r)
			assert.Equal(t, returnedError, addAssuranceVisitReturnedError)
		})
	}
}
