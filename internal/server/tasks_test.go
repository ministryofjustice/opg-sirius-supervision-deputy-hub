package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockTasksClient struct {
	count       int
	lastCtx     sirius.Context
	err         error
	taskTypes   []model.TaskType
	teamMembers []model.TeamMember
	tasks       sirius.TaskList
}

func (m *mockTasksClient) GetTaskTypesForDeputyType(ctx sirius.Context, deputyId string) ([]model.TaskType, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskTypes, m.err
}

func (m *mockTasksClient) GetDeputyTeamMembers(ctx sirius.Context, defaultPATeam int, deputy sirius.DeputyDetails) ([]model.TeamMember, error) {
	m.count += 1

	return m.teamMembers, m.err
}

func (m *mockTasksClient) GetTasks(ctx sirius.Context, deputyId int) (sirius.TaskList, error) {
	m.count += 1

	return m.tasks, m.err
}

func TestNavigateTasksTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockTasksClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForTasks(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
