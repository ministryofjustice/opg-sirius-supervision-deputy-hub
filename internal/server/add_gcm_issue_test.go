package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddGCMIssueClient struct {
	count          int
	AddGCMIssueErr error
	GcmIssueTypes  []model.RefData
	CaseRecNumber  string
	Client         sirius.DeputyClient
	HasFoundClient string
	GcmIssueType   model.RefData
	Notes          string
}

func (m *mockAddGCMIssueClient) GetGCMIssueTypes(ctx sirius.Context) ([]model.RefData, error) {
	return []model.RefData{
		{
			"MISSING_INFORMATION",
			"Missing information",
			false,
		},
	}, nil
}

func (m *mockAddGCMIssueClient) GetDeputyClient(ctx sirius.Context, caseRecNumber string, deputyId int) (sirius.DeputyClient, error) {
	return sirius.DeputyClient{
		ClientId:  1234,
		Firstname: "Test",
		CourtRef:  "123456",
	}, nil
}

func (m *mockAddGCMIssueClient) AddGcmIssue(ctx sirius.Context, caseRecNumber, notes string, gcmIssueType model.RefData, deputyId int) error {
	return nil
}

var addGCMIssueAppVars = AppVars{
	DeputyDetails: testDeputy,
}

func TestGetAddGCMIssue(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddGCMIssueClient{}
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

func TestAddGCMIssueCanSearchForClient(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddGCMIssueClient{}
	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
	}
	template := &mockTemplates{}

	form := url.Values{
		"case-number":       {"123456"},
		"search-for-client": {"search-for-client"},
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
				"MISSING_INFORMATION",
				"Missing information",
				false,
			},
		},
		CaseRecNumber: "123456",
		Client: sirius.DeputyClient{
			ClientId:  1234,
			CourtRef:  "123456",
			Firstname: "Test",
		},
		HasFoundClient: "",
		GcmIssueType:   model.RefData{},
		Notes:          "",
	}, template.lastVars)

	assert.Nil(res)
}

func TestAddGCMIssueCanAddGcmIssue(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddGCMIssueClient{
		Client: sirius.DeputyClient{
			ClientId:  1234,
			CourtRef:  "123456",
			Firstname: "Test",
		},
	}
	app := AppVars{
		Path:          "/path",
		DeputyDetails: sirius.DeputyDetails{ID: 123},
	}
	template := &mockTemplates{}

	form := url.Values{
		"case-number":       {"123456"},
		"search-for-client": {""},
		"issue-type":        {"Missing information"},
		"notes":             {"test note"},
		"submit-form":       {"submit-form"},
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
	assert.Equal(res, Redirect("/123/gcm-issues?success=addGcmIssue&123456"))

}
