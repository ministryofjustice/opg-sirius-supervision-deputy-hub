package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddAssuranceVisitInformation struct {
	count   int
	lastCtx sirius.Context
	err     error
	userDetails sirius.UserDetails
}

func (m *mockAddAssuranceVisitInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetails, m.err
}

func (m *mockAddAssuranceVisitInformation) AddAssuranceVisit(ctx sirius.Context, requestedDate string, userId, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestPostAssuranceVisit(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddAssuranceVisitInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddAssuranceVisit(client, nil)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/123/assurance-visits?success=addAssuranceVisit"))
}

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
//	client.err = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//	defaultPATeam := 23
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
//	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForAddFirm(client, defaultPATeam, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//
//	assert.Equal(addFirmVars{
//		Path:   "/133",
//		Errors: validationErrors,
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}


