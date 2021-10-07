package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
)

type mockEditDeputyHubInformation struct {
	count      int
	lastCtx    sirius.Context
	err        error
	deputyData sirius.DeputyDetails
	editDeputyData sirius.DeputyDetails
}

func (m *mockEditDeputyHubInformation) GetDeputyDetails(ctx sirius.Context, deputyId int) (sirius.DeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockEditDeputyHubInformation) EditDeputyDetails(ctx sirius.Context, deputyDetails sirius.DeputyDetails)  error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestNavigateToEditDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditDeputyHubInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestErrorEditDeputyMessageWhenStringLengthTooLong(t *testing.T) {
	assert := assert.New(t)
	client := &mockEditDeputyHubInformation{}

	validationErrors := sirius.ValidationErrors{
		"organisationName": {
			"stringLengthTooLong": "What sirius gives us",
		},"workPhoneNumber": {
			"stringLengthTooLong": "What sirius gives us",
		},"email": {
			"stringLengthTooLong": "What sirius gives us",
		},"organisationTeamOrDepartmentName": {
			"stringLengthTooLong": "What sirius gives us",
		},"addressLine1": {
			"stringLengthTooLong": "What sirius gives us",
		},"addressLine2": {
			"stringLengthTooLong": "What sirius gives us",
		},"addressLine3": {
			"stringLengthTooLong": "What sirius gives us",
		},"town": {
			"stringLengthTooLong": "What sirius gives us",
		},"county": {
			"stringLengthTooLong": "What sirius gives us",
		},"postcode": {
			"stringLengthTooLong": "What sirius gives us",
		},
	}

	client.err = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter();
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForEditDeputyHub(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"organisationName": {
            "stringLengthTooLong": "The deputy name must be 255 characters or fewer",
        },"workPhoneNumber": {
            "stringLengthTooLong": "The telephone number must be 255 characters or fewer",
        },"email": {
            "stringLengthTooLong": "The email number must be 255 characters or fewer",
        },"organisationTeamOrDepartmentName": {
            "stringLengthTooLong": "The team or department must be 255 characters or fewer",
        },"addressLine1": {
            "stringLengthTooLong": "The building or street must be 255 characters or fewer",
        },"addressLine2": {
            "stringLengthTooLong": "Address line 2 must be 255 characters or fewer",
        },"addressLine3": {
            "stringLengthTooLong": "AddressLine 3 must be 255 characters or fewer",
        },"town": {
            "stringLengthTooLong": "The town or city must be 255 characters or fewer",
        },"county": {
            "stringLengthTooLong": "The county must be 255 characters or fewer",
        },"postcode": {
            "stringLengthTooLong": "The postcode must be 255 characters or fewer",
        },
	}

	assert.Equal(2, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editDeputyHubVars{
		Path:    "/123",
		Errors:  expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestErrorEditDeputyMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockEditDeputyHubInformation{}

	validationErrors := sirius.ValidationErrors{
		"organisationName": {
			"isEmpty": "What sirius gives us",
		},
	}

	client.err = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter();
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForEditDeputyHub(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"organisationName": {
            "isEmpty": "Enter a deputy name",
        },
	}

	assert.Equal(2, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editDeputyHubVars{
		Path:    "/123",
		Errors:  expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
