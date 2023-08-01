package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManageTasks struct {
	mock.Mock
}

func (m *mockManageTasks) GetTask(ctx sirius.Context, taskId int) (model.Task, error) {
	args := m.Called(ctx)

	return args.Get(0).(model.Task), args.Error(1)
}

func (m *mockManageTasks) GetDeputyTeamMembers(ctx sirius.Context, defaultPATeam int, deputy sirius.DeputyDetails) ([]model.TeamMember, error) {
	args := m.Called(ctx)

	return args.Get(0).([]model.TeamMember), args.Error(1)
}

func (m *mockManageTasks) UpdateTask(ctx sirius.Context, deputyId, taskId int, dueDate, notes string, assigneeId int) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *mockManageTasks) GetTaskTypesForDeputyType(ctx sirius.Context, deputyType string) ([]model.TaskType, error) {

	args := m.Called(ctx)

	return args.Get(0).([]model.TaskType), args.Error(1)
}

func TestNavigateToManageTask(t *testing.T) {
	assert := assert.New(t)

	deputyDetails := sirius.DeputyDetails{ID: 123}
	task := model.Task{Id: 555}
	teamMembers := []model.TeamMember{{ID: 99}}
	taskTypes := []model.TaskType{{Handle: "ABC", Description: "A Big Critical Task"}}

	client := &mockManageTasks{}
	client.On("GetTask", mock.Anything).Return(task, nil)
	client.On("GetDeputyTeamMembers", mock.Anything).Return(teamMembers, nil)
	client.On("GetTaskTypesForDeputyType", mock.Anything).Return(taskTypes, nil)

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageTasks(client, template)
	err := handler(sirius.DeputyDetails{ID: 123}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(manageTaskVars{
		DeputyDetails: deputyDetails,
		TaskDetails:   task,
		Assignees:     teamMembers,
	}, template.lastVars)
}

func TestPostManageTask(t *testing.T) {
	assert := assert.New(t)

	deputyDetails := sirius.DeputyDetails{ID: 123}
	task := model.Task{Id: 555}
	teamMembers := []model.TeamMember{{ID: 99}}
	taskTypes := []model.TaskType{{Handle: "ABC", Description: "A Big Critical Task"}}

	client := &mockManageTasks{}
	client.On("GetTask", mock.Anything).Return(task, nil)
	client.On("GetDeputyTeamMembers", mock.Anything).Return(teamMembers, nil)
	client.On("GetTaskTypesForDeputyType", mock.Anything).Return(taskTypes, nil)
	client.On("UpdateTask", mock.Anything).Return(nil)

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/tasks/123", strings.NewReader("{dueDate:'2200/10/20'}"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForManageTasks(client, template)(deputyDetails, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(Redirect("/123/tasks?success=manageTask"), redirect)
}
