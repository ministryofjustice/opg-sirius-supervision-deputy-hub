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
	deputyEvents       sirius.TimelineList
	pagination         paginate.Pagination
}

func (m *mockDeputyHubTimelineInformation) GetDeputyEvents(ctx sirius.Context, deputyId int, pageNumber int, timelineEventsPerPage int) (sirius.TimelineList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyEvents, m.GetDeputyEventsErr
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
	assert.Equal(paginate.Pagination{
		CurrentPage:     0,
		TotalPages:      0,
		TotalElements:   0,
		ElementsPerPage: 0,
		ElementName:     "",
		PerPageOptions:  nil,
		UrlBuilder:      nil,
	}, client.pagination)
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
