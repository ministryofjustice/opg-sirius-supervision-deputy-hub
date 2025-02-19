package server

import (
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

	handler := renderTemplateForDeputyHubEvents(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestDeputyEventsReturnsErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyHubTimelineInformation{
		GetDeputyEventsErr: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", strings.NewReader(""))

	returnedError := renderTemplateForDeputyHubEvents(client, template)(AppVars{}, w, r)

	assert.Equal(client.GetDeputyEventsErr, returnedError)
}
