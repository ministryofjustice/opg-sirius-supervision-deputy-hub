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

type mockAddAssuranceVisitInformation struct {
	count   int
	lastCtx sirius.Context
	err     error
}

func (m *mockAddAssuranceVisitInformation) AddAssuranceVisit(ctx sirius.Context, assuranceType string, requestedDate string, userId, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestPostAssuranceVisit(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddAssuranceVisitInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurance-visits", strings.NewReader("{requestedDate:'2200/10/20', requestedBy:22}"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/assurance-visits", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddAssuranceVisit(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Nil(returnedError)
}
