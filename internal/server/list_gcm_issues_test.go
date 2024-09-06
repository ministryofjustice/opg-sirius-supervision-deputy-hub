package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockGcmIssues struct {
	mock.Mock
}

func (m *mockGcmIssues) CloseGCMIssues(ctx sirius.Context, gcmIds []string) error {
	args := m.Called(ctx, gcmIds)

	return args.Error(0)
}

func (m *mockGcmIssues) GetGCMIssues(ctx sirius.Context, deputyId int, params sirius.GcmIssuesParams) ([]sirius.GcmIssue, error) {
	args := m.Called(ctx, deputyId)

	return args.Get(0).([]sirius.GcmIssue), args.Error(1)
}

func TestNavigateToGcmIssuesTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockGcmIssues{}

	deputyDetails := sirius.DeputyDetails{ID: 123}
	app := AppVars{
		DeputyDetails: deputyDetails,
		PageName:      "General Case Manager issues",
	}
	gcmIssues := []sirius.GcmIssue{
		{
			Id:            1,
			Client:        sirius.GcmClient{},
			CreatedDate:   "2024-01-01",
			CreatedByUser: sirius.UserInformation{},
			Notes:         "Problem here",
			GcmIssueType: model.RefData{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
		},
	}

	client.On("GetGCMIssues", mock.Anything, 123).Return(gcmIssues, nil)

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForGcmIssues(client, template)
	err := handler(app, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(GcmIssuesVars{
		AppVars:   app,
		GcmIssues: gcmIssues,
		Sort: urlbuilder.Sort{
			OrderBy:    "createdDate",
			Descending: false,
		},
		UrlBuilder: urlbuilder.UrlBuilder{
			OriginalPath: "open-issues",
			SelectedSort: urlbuilder.Sort{
				OrderBy:    "createdDate",
				Descending: false,
			},
		},
	}, template.lastVars)
}

func TestCloseGCMIssue(t *testing.T) {
	assert := assert.New(t)
	client := &mockGcmIssues{}
	deputyDetails := sirius.DeputyDetails{ID: 123}
	app := AppVars{
		DeputyDetails: deputyDetails,
		PageName:      "General Case Manager issues",
	}

	client.On("CloseGCMIssues", mock.Anything, []string{"25"}).Return(nil)

	template := &mockTemplates{}

	form := url.Values{"selected-gcms": []string{"25"}}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/gcm-issues/closed-issues", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/gcm-issues/closed-issues", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForGcmIssues(client, template)(app, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(Redirect("/123/gcm-issues/open-issues?success=closedGcms&count=1"), redirect)
}
