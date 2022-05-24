package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManageDeputyContactDetailsInformation struct {
	count      int
	lastCtx    sirius.Context
	err        error
	deputyData sirius.DeputyDetails
	updateErr  error
}

func (m *mockManageDeputyContactDetailsInformation) UpdateDeputyContactDetails(ctx sirius.Context, _ int, _ sirius.DeputyContactDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.updateErr
}

func TestGetManageDeputyDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyContactDetailsInformation{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageDeputyContactDetails(client, defaultPATeam, template)
	err := handler(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPostManageDeputyDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyContactDetailsInformation{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForManageDeputyContactDetails(client, defaultPATeam, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(redirect, Redirect("/123?success=deputyDetails"))
}

func TestErrorManageDeputyDetailsMessageWhenStringLengthTooLong(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageDeputyContactDetailsInformation{}

	validationErrors := sirius.ValidationErrors{
		"firstname": {
			"stringLengthTooLong": "What sirius gives us",
		}, "surname": {
			"stringLengthTooLong": "What sirius gives us",
		}, "organisationName": {
			"stringLengthTooLong": "What sirius gives us",
		}, "workPhoneNumber": {
			"stringLengthTooLong": "What sirius gives us",
		}, "email": {
			"stringLengthTooLong": "What sirius gives us",
		}, "addressLine1": {
			"stringLengthTooLong": "What sirius gives us",
		}, "addressLine2": {
			"stringLengthTooLong": "What sirius gives us",
		}, "addressLine3": {
			"stringLengthTooLong": "What sirius gives us",
		}, "town": {
			"stringLengthTooLong": "What sirius gives us",
		}, "county": {
			"stringLengthTooLong": "What sirius gives us",
		}, "postcode": {
			"stringLengthTooLong": "What sirius gives us",
		},
	}

	client.updateErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageDeputyContactDetails(client, defaultPATeam, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"firstname": {
			"stringLengthTooLong": "The deputy first name must be 255 characters or fewer",
		},
		"surname": {
			"stringLengthTooLong": "The deputy surname must be 255 characters or fewer",
		},
		"organisationName": {
			"stringLengthTooLong": "The organisation name must be 255 characters or fewer",
		},
		"workPhoneNumber": {
			"stringLengthTooLong": "The telephone number must be 255 characters or fewer",
		},
		"email": {
			"stringLengthTooLong": "The email must be 255 characters or fewer",
		},
		"addressLine1": {
			"stringLengthTooLong": "The building or street must be 255 characters or fewer",
		},
		"addressLine2": {
			"stringLengthTooLong": "Address line 2 must be 255 characters or fewer",
		},
		"addressLine3": {
			"stringLengthTooLong": "Address line 3 must be 255 characters or fewer",
		},
		"town": {
			"stringLengthTooLong": "The town or city must be 255 characters or fewer",
		},
		"county": {
			"stringLengthTooLong": "The county must be 255 characters or fewer",
		},
		"postcode": {
			"stringLengthTooLong": "The postcode must be 255 characters or fewer",
		},
	}

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(manageDeputyContactDetailsVars{
		Path:     "/123",
		DeputyId: 123,
		Errors:   expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestErrorManageDeputyDetailsMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageDeputyContactDetailsInformation{}

	validationErrors := sirius.ValidationErrors{
		"firstname": {
			"isEmpty": "What sirius gives us",
		},
		"surname": {
			"isEmpty": "What sirius gives us",
		},
		"organisationName": {
			"isEmpty": "What sirius gives us",
		},
	}

	client.updateErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageDeputyContactDetails(client, defaultPATeam, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"firstname": {
			"isEmpty": "The deputy first name is required and can't be empty",
		},
		"surname": {
			"isEmpty": "The deputy surname is required and can't be empty",
		},
		"organisationName": {
			"isEmpty": "The organisation name is required and can't be empty",
		},
	}

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(manageDeputyContactDetailsVars{
		Path:     "/123",
		DeputyId: 123,
		Errors:   expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
