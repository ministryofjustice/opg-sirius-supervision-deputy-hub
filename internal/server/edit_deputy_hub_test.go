package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
// 	"github.com/gorilla/mux"
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

func (m *mockEditDeputyHubInformation) EditDeputyDetails(ctx sirius.Context, deputyId int) (sirius.DeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.editDeputyData, m.err
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

// func TestErrorEditDeputyMessageWhenStringLengthTooLong(t *testing.T) {
// 	assert := assert.New(t)
// 	client := &mockEditDeputyHubInformation{}
//
// 	validationErrors := sirius.ValidationErrors{
// 		"organisationName": {
// 			"stringLengthTooLong": "The deputy name must be 255 characters or fewer",
// 		},
// 		"description": {
// 			"stringLengthTooLong": "The deputy name must be 255 characters or fewer",
// 		},
// 	}
// 	client.err = sirius.ValidationError{
// 		Errors: validationErrors,
// 	}
//
// 	template := &mockTemplates{}
//
// 	w := httptest.NewRecorder()
// 	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
// 	var returnedError error
//
// 	testHandler := mux.NewRouter();
// 	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
// 		returnedError = renderTemplateForEditDeputyHub(client, template)(sirius.PermissionSet{}, w, r)
// 	})
//
// 	testHandler.ServeHTTP(w, r)
//
// 	expectedValidationErrors := sirius.ValidationErrors{
// 		"1-title": {
// 			"stringLengthTooLong": "The title must be 255 characters or fewer",
// 		},
// 		"2-note": {
// 			"stringLengthTooLong": "The note must be 1000 characters or fewer",
// 		},
// 	}
//
// 	assert.Equal(3, client.count)
//
// 	assert.Equal(1, template.count)
// 	assert.Equal("page", template.lastName)
// 	assert.Equal(editDeputyHubVars{
// 		Path:    "/123",
// 		Errors:  expectedValidationErrors,
// 	}, template.lastVars)
//
// 	assert.Nil(returnedError)
// }
