package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddTasksClient struct {
	mock.Mock
	count                   int
	lastCtx                 sirius.Context
	AddTaskErr              error
	GetTaskTypesErr         error
	GetDeputyTeamMembersErr error
	taskTypes               []model.TaskType
	assignees               []model.TeamMember
	selectedAssignee        int
}

func (m *mockAddTasksClient) AddTask(ctx sirius.Context, deputyId int, taskType string, typeName string, dueDate string, notes string, assigneeId int) error {
	m.count += 1
	m.lastCtx = ctx
	m.selectedAssignee = assigneeId

	return m.AddTaskErr
}

func (m *mockAddTasksClient) GetTaskTypesForDeputyType(ctx sirius.Context, details string) ([]model.TaskType, error) {

	m.count += 1
	m.lastCtx = ctx

	return m.taskTypes, m.GetTaskTypesErr
}

func (m *mockAddTasksClient) GetDeputyTeamMembers(ctx sirius.Context, defaultPaTeam int, details sirius.DeputyDetails) ([]model.TeamMember, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.assignees, m.GetDeputyTeamMembersErr
}

func TestGetTasks(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTasksClient{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddTask(client, defaultPATeam, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestLoadAddTaskForm(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTasksClient{}
	template := &mockTemplates{}
	defaultPATeam := 23

	deputy := sirius.DeputyDetails{ID: 1, ExecutiveCaseManager: sirius.ExecutiveCaseManager{
		EcmId: 1,
	}}
	taskTypes := []model.TaskType{{Handle: "ABC"}}
	client.taskTypes = taskTypes
	assignees := []model.TeamMember{{ID: 1, DisplayName: "Teamster"}}
	client.assignees = assignees

	expectedVars := AddTaskVars{
		Path:          "/path",
		DeputyDetails: deputy,
		TaskTypes:     taskTypes,
		Assignees:     assignees,
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddTask(client, defaultPATeam, template)
	res := handler(deputy, w, r)

	assert.Nil(res)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(expectedVars, template.lastVars)
}

func TestAddTask_success_ecm(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddTasksClient{}
	defaultPATeam := 23

	deputy := sirius.DeputyDetails{ID: 123}
	taskTypes := []model.TaskType{{Handle: "ABC", Description: "A Big Critical Task"}}
	client.taskTypes = taskTypes
	assignees := []model.TeamMember{{ID: 1, DisplayName: "Teamster"}}
	client.assignees = assignees

	form := url.Values{
		"tasktype":          {"ABC"},
		"duedate":           {"2022-04-02"},
		"notes":             {"A note"},
		"assignedto":        {"1"},
		"select-assignedto": {},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var res error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		res = renderTemplateForAddTask(client, defaultPATeam, nil)(deputy, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(res, Redirect("/123/tasks?success=add&taskType=A Big Critical Task"))
	assert.Equal(1, client.selectedAssignee)
}

func TestAddTask_success_other(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddTasksClient{}
	defaultPATeam := 23

	deputy := sirius.DeputyDetails{ID: 123}
	taskTypes := []model.TaskType{{Handle: "ABC", Description: "A Big Critical Task"}}
	client.taskTypes = taskTypes
	assignees := []model.TeamMember{{ID: 1, DisplayName: "Teamster"}}
	client.assignees = assignees

	form := url.Values{
		"tasktype":          {"ABC"},
		"duedate":           {"2022-04-02"},
		"notes":             {"A note"},
		"assignedto":        {"other"},
		"select-assignedto": {"2"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var res error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		res = renderTemplateForAddTask(client, defaultPATeam, nil)(deputy, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(res, Redirect("/123/tasks?success=add&taskType=A Big Critical Task"))
	assert.Equal(2, client.selectedAssignee)
}

func TestAddTaskValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddTasksClient{}
	defaultPATeam := 23

	validationErrors := sirius.ValidationErrors{
		"dueDate": {
			"dateFalseFormat": "This must be a real date",
		},
	}

	client.AddTaskErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res := renderTemplateForAddTask(client, defaultPATeam, template)(sirius.DeputyDetails{}, w, r)

	assert.Equal(AddTaskVars{
		Path:   "/133",
		Errors: validationErrors,
	}, template.lastVars)

	assert.Nil(res)
}

func TestAddTasksHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockAddTasksClient
	}{
		{
			Client: &mockAddTasksClient{
				GetTaskTypesErr: returnedError,
			},
		},
		{
			Client: &mockAddTasksClient{
				GetDeputyTeamMembersErr: returnedError,
			},
		},
		{
			Client: &mockAddTasksClient{
				AddTaskErr: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))

			addTaskReturnedError := renderTemplateForAddTask(client, 23, template)(sirius.DeputyDetails{}, w, r)
			assert.Equal(t, returnedError, addTaskReturnedError)
		})
	}
}