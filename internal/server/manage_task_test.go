package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManageTasks struct {
	count          int
	lastCtx        sirius.Context
	err            error
	DeputyDetails  sirius.DeputyDetails
	TaskDetails    model.Task
	Error          string
	Errors         sirius.ValidationErrors
	Success        bool
	SuccessMessage string
	Assignees      []model.TeamMember
}

func (m *mockManageTasks) GetTask(ctx sirius.Context, taskId int) (model.Task, error) {
	m.count += 1
	m.lastCtx = ctx
	return m.TaskDetails, m.err
}

func (m *mockManageTasks) GetDeputyTeamMembers(ctx sirius.Context, defaultPATeam int, deputy sirius.DeputyDetails) ([]model.TeamMember, error) {
	m.count += 1
	m.lastCtx = ctx
	return m.Assignees, m.err
}

func (m *mockManageTasks) UpdateTask(ctx sirius.Context, deputyId, taskId int, dueDate, notes string, assigneeId int) error {
	m.count += 1
	m.lastCtx = ctx
	return m.err
}

func (m *mockManageTasks) GetTaskTypesForDeputyType(ctx sirius.Context, deputyType string) ([]model.TaskType, error) {

	m.count += 1
	m.lastCtx = ctx
	return nil, m.err
}

func TestNavigateToChangeTask(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageTasks{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForManageTasks(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
