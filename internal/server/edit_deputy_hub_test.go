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

type mockEditDeputyHubInformation struct {
	count            int
	lastCtx          sirius.Context
	err              error
	deputyClientData sirius.ClientList
	ariaSorting      sirius.AriaSorting
}

func (m *mockEditDeputyHubInformation) EditDeputyDetails(ctx sirius.Context, deputyDetails sirius.DeputyDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockEditDeputyHubInformation) GetDeputyClients(ctx sirius.Context, deputyId, displayClientLimit, search int, deputyType, columnBeingSorted, sortOrder string) (sirius.ClientList, sirius.AriaSorting, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.ariaSorting, m.err
}

func TestNavigateToEditDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditDeputyHubInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestEditDeputyValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockEditDeputyHubInformation{}

	validationErrors := sirius.ValidationErrors{
		"organisationName": {
			"stringLengthTooLong": "What sirius gives us",
		},
	}

	client.err = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForEditDeputyHub(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editDeputyHubVars{
		AppVars: AppVars{
			Errors: validationErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}
