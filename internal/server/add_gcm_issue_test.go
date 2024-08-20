package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddGCMIssueClient struct {
	GetGCMIssueTypesError error
	GetDeputyClientError  error
	AddGCMIssueErr        error
	GcmIssueTypes         []model.RefData
	CaseRecNumber         string
	Client                sirius.DeputyClient
	GcmIssueType          string
	Notes                 string
}

func (m *mockAddGCMIssueClient) GetGCMIssueTypes(ctx sirius.Context) ([]model.RefData, error) {
	return m.GcmIssueTypes, m.GetGCMIssueTypesError
}

func (m *mockAddGCMIssueClient) GetDeputyClient(ctx sirius.Context, caseRecNumber string, deputyId int) (sirius.DeputyClient, error) {
	return m.Client, m.GetDeputyClientError
}

func (m *mockAddGCMIssueClient) AddGcmIssue(ctx sirius.Context, caseRecNumber, notes string, gcmIssueType string, deputyId int) error {
	return m.AddGCMIssueErr
}

var addGCMIssueAppVars = AppVars{
	DeputyDetails: testDeputy,
}

func TestGetAddGCMIssue(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddGCMIssueClient{
		GcmIssueTypes: []model.RefData{
			{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
		},
	}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddGcmIssue(client, template)
	err := handler(addGCMIssueAppVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostAddGCMIssueSearchForClient(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddGCMIssueClient{
		GcmIssueTypes: []model.RefData{
			{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
		},
		Client: sirius.DeputyClient{
			ClientId:  1234,
			Firstname: "Test",
			CourtRef:  "123456",
		},
	}
	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
	}
	template := &mockTemplates{}

	form := url.Values{
		"case-number":       {"123456"},
		"search-for-client": {"true"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var res error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		res = renderTemplateForAddGcmIssue(client, template)(app, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(AddGcmIssueVars{
		AppVars: AppVars{
			Path:          "/path",
			PageName:      "Add a GCM issue",
			DeputyDetails: sirius.DeputyDetails{ID: 123},
		},
		GcmIssueTypes: []model.RefData{
			{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
		},
		CaseRecNumber: "123456",
		Client: sirius.DeputyClient{
			ClientId:  1234,
			CourtRef:  "123456",
			Firstname: "Test",
		},
		GcmIssueType: "",
		Notes:        "",
	}, template.lastVars)

	assert.Nil(res)
}

func TestPostAddGCMIssueSubmitForm(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddGCMIssueClient{
		Client: sirius.DeputyClient{
			ClientId: 1234,
		},
		GcmIssueTypes: []model.RefData{
			{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
		},
	}
	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
	}
	template := &mockTemplates{}

	form := url.Values{
		"case-number": {"123456"},
		"submit-form": {"true"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var res error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		res = renderTemplateForAddGcmIssue(client, template)(app, w, r)
	})

	testHandler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(res, Redirect("/123/gcm-issues/open-issues?success=addGcmIssue&123456"))
}

func TestPostAddGCMIssueErrorsIfNoCaseNumber(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddGCMIssueClient{
		Client: sirius.DeputyClient{
			ClientId: 1234,
		},
		GcmIssueTypes: []model.RefData{
			{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
		},
	}

	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
		PageName:      "Add a GCM issue",
		Errors: sirius.ValidationErrors{
			"client-case-number": {
				"": "Enter a case number",
			},
		},
	}
	template := &mockTemplates{}

	form := url.Values{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var res error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		res = renderTemplateForAddGcmIssue(client, template)(app, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Nil(res)

	assert.Equal(AddGcmIssueVars{
		AppVars: app,
		GcmIssueTypes: []model.RefData{
			{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
		},
	}, template.lastVars)
}

func TestPostAddGCMIssueHandlesErrorsInOtherClientFiles(t *testing.T) {
	expectedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockAddGCMIssueClient
		Form   url.Values
	}{
		{
			Client: &mockAddGCMIssueClient{
				GetGCMIssueTypesError: expectedError,
			},
			Form: url.Values{
				"case-number":       {"123456"},
				"search-for-client": {"true"},
			},
		},
		{
			Client: &mockAddGCMIssueClient{
				GetDeputyClientError: expectedError,
			},
			Form: url.Values{
				"case-number":       {"123456"},
				"search-for-client": {"true"},
			},
		},
		{
			Client: &mockAddGCMIssueClient{
				AddGCMIssueErr: expectedError,
				Client:         sirius.DeputyClient{ClientId: 1},
			},
			Form: url.Values{
				"case-number": {"123456"},
				"submit-form": {"true"},
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {
			app := AppVars{
				Path:          "/path",
				DeputyDetails: sirius.DeputyDetails{ID: 123},
			}
			template := &mockTemplates{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/123", strings.NewReader(tc.Form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			var res error

			testHandler := mux.NewRouter()
			testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
				res = renderTemplateForAddGcmIssue(tc.Client, template)(app, w, r)
			})

			testHandler.ServeHTTP(w, r)

			assert.Equal(t, expectedError, res)
		})
	}
}
