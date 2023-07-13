package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubTimelineInformation struct {
	count        int
	lastCtx      sirius.Context
	err          error
	deputyEvents sirius.DeputyEvents
}

func (m *mockDeputyHubTimelineInformation) GetDeputyEvents(ctx sirius.Context, deputyId int) (sirius.DeputyEvents, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyEvents, m.err
}

func TestNavigateToTimeline(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubTimelineInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHubEvents(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
