package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubTimelineInformation struct {
	count              int
	lastCtx            sirius.Context
	GetDeputyEventsErr error
}

func (m *mockDeputyHubTimelineInformation) GetDeputyEvents(ctx sirius.Context, deputyId int, pageNumber int, timelineEventsPerPage int) (sirius.TimelineList, error) {
	m.count += 1
	m.lastCtx = ctx

	var getDeputyEventsVars = sirius.TimelineList{
		Limit: 25,
		Pages: struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		}{
			Current: 1,
			Total:   1,
		},
		Total: 2,
		DeputyEvents: sirius.DeputyEvents{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		},
	}

	return getDeputyEventsVars, m.GetDeputyEventsErr
}

func TestNavigateToTimeline(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubTimelineInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHubEvents(client, template, EnvironmentVars{})
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(deputyHubEventVars{
		AppVars: AppVars{
			PageName: "Timeline",
		},
		DeputyEvents: sirius.DeputyEvents{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		},
		Pagination: paginate.Pagination{
			CurrentPage:     1,
			TotalPages:      1,
			TotalElements:   2,
			ElementsPerPage: 25,
			ElementName:     "timeline event(s)",
			PerPageOptions:  []int{25, 50, 100},
			UrlBuilder: urlbuilder.UrlBuilder{
				OriginalPath:    "/path",
				SelectedPerPage: 25,
			},
		},
	}, template.lastVars)
}

func TestDeputyEventsReturnsErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyHubTimelineInformation{
		GetDeputyEventsErr: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", strings.NewReader(""))

	returnedError := renderTemplateForDeputyHubEvents(client, template, EnvironmentVars{})(AppVars{}, w, r)

	assert.Equal(client.GetDeputyEventsErr, returnedError)
}

func TestCreateUrlBuilder(t *testing.T) {
	myUrlBuilder := CreateUrlBuilder("67/timeline", 25, "/supervision/deputies/")

	assert.Equal(t, urlbuilder.UrlBuilder{
		OriginalPath:    "/supervision/deputies/67/timeline",
		SelectedPerPage: 25,
		SelectedFilters: nil,
		SelectedSort:    urlbuilder.Sort{},
	}, myUrlBuilder)

	myUrlBuilder = CreateUrlBuilder("99/url", 100, "/test/path/")

	assert.Equal(t, urlbuilder.UrlBuilder{
		OriginalPath:    "/test/path/99/url",
		SelectedPerPage: 100,
		SelectedFilters: nil,
		SelectedSort:    urlbuilder.Sort{},
	}, myUrlBuilder)
}
