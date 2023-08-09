package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyChangeFirmInformation struct {
	count                 int
	lastCtx               sirius.Context
	GetFirmsErr           error
	AssignDeputyToFirmErr error
	firmData              []sirius.FirmForList
}

func (m *mockDeputyChangeFirmInformation) GetFirms(ctx sirius.Context) ([]sirius.FirmForList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firmData, m.GetFirmsErr
}

func (m *mockDeputyChangeFirmInformation) AssignDeputyToFirm(ctx sirius.Context, defaultPATeam int, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.AssignDeputyToFirmErr
}

func TestNavigateToChangeFirm(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyChangeFirmInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForChangeFirm(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPostChangeFirm(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyChangeFirmInformation{}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/firm", strings.NewReader(""))

	form := url.Values{}
	form.Add("select-existing-firm", "26")
	r.PostForm = form

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/firm", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForChangeFirm(client, template)(sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(returnedError, Redirect("/76?success=firm"))
}

func TestPostChangeFirmReturnsStatusMethodError(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyChangeFirmInformation{}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/76/firm", strings.NewReader(""))

	returnedError := renderTemplateForChangeFirm(client, template)(sirius.DeputyDetails{}, w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), returnedError)
}

func TestPostFirmRedirectsIfNewFirmChosen(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyChangeFirmInformation{}

	template := &mockTemplates{}

	client.AssignDeputyToFirmErr = StatusError(http.StatusGatewayTimeout)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/firm", strings.NewReader(""))
	form := url.Values{}
	form.Add("select-firm", "new-firm")
	r.PostForm = form

	returnedError := renderTemplateForChangeFirm(client, template)(sirius.DeputyDetails{}, w, r)

	assert.Equal(Redirect("/0/add-firm"), returnedError)
}

func TestPostFirmReturnsValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyChangeFirmInformation{}

	validationErrors := sirius.ValidationErrors{
		"changeFirm": {
			"stringLengthTooLong": "The firm name must be 255 characters or fewer",
		},
	}

	client.AssignDeputyToFirmErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/firm", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error
	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/firm", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForChangeFirm(client, template)(sirius.DeputyDetails{}, w, r)
	})
	testHandler.ServeHTTP(w, r)

	assert.Equal(changeFirmVars{
		Path:   "/76/firm",
		Errors: validationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)

}

func TestChangeFirmHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockDeputyChangeFirmInformation
	}{
		{
			Client: &mockDeputyChangeFirmInformation{
				GetFirmsErr: returnedError,
			},
		},
		{
			Client: &mockDeputyChangeFirmInformation{
				AssignDeputyToFirmErr: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {

			client := tc.Client
			form := url.Values{}
			form.Add("select-firm", "")
			form.Add("select-existing-firm", "26")

			template := &mockTemplates{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/76/firm", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.PostForm = form

			changeEcmReturnedError := renderTemplateForChangeFirm(client, template)(sirius.DeputyDetails{}, w, r)

			assert.Equal(t, returnedError, changeEcmReturnedError)
		})
	}
}
