package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockCompleteTaskClient struct {
	mock.Mock
	count                        int
	lastCtx                      sirius.Context
	GetTaskTypesForDeputyTypeErr error
	CompleteTaskErr              error
	GetTaskErr                   error
	taskTypes                    []model.TaskType
	task                         model.Task
}

func (m *mockCompleteTaskClient) GetTaskTypesForDeputyType(ctx sirius.Context, details string) ([]model.TaskType, error) {

	m.count += 1
	m.lastCtx = ctx

	return m.taskTypes, m.GetTaskTypesForDeputyTypeErr
}

func (m *mockCompleteTaskClient) CompleteTask(ctx sirius.Context, taskId int, notes string) error {
	m.count += 1
	m.lastCtx = ctx

	return m.CompleteTaskErr
}

func (m *mockCompleteTaskClient) GetTask(ctx sirius.Context, deputyId int) (model.Task, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.task, m.GetTaskErr
}

func TestGetCompleteTask(t *testing.T) {
	assert := assert.New(t)

	client := &mockCompleteTaskClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForCompleteTask(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestLoadCompleteTaskForm(t *testing.T) {
	assert := assert.New(t)

	deputy := sirius.DeputyDetails{ID: 1, ExecutiveCaseManager: sirius.ExecutiveCaseManager{
		EcmId: 1,
	}}
	taskDetails := model.Task{}

	client := &mockCompleteTaskClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForCompleteTask(client, template)
	res := handler(deputy, w, r)

	assert.Nil(res)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(completeTaskVars{
		Path:          "/path",
		DeputyDetails: deputy,
		TaskDetails:   taskDetails,
	}, template.lastVars)
}

func TestPostCompleteTask(t *testing.T) {
	assert := assert.New(t)

	deputy := sirius.DeputyDetails{ID: 123}
	taskDetails := model.Task{Type: "ABC"}
	taskTypes := []model.TaskType{{Handle: "ABC", Description: "TaskDescription"}}

	client := &mockCompleteTaskClient{}
	client.taskTypes = taskTypes
	client.task = taskDetails

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	form := url.Values{"taskCompletedNotes": {"some notes"}}
	r, _ := http.NewRequest("POST", "/path", strings.NewReader(form.Encode()))

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForCompleteTask(client, template)(deputy, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(Redirect("/123/tasks?success=complete&taskType=TaskDescription"), redirect)
}

func TestCompleteTaskValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockCompleteTaskClient{}

	validationErrors := sirius.ValidationErrors{
		"taskCompletedNotes": {
			"stringLengthTooLong": "This value is too long",
		},
	}

	client.CompleteTaskErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/111", strings.NewReader(""))
	returnedError := renderTemplateForCompleteTask(client, template)(sirius.DeputyDetails{}, w, r)

	assert.Equal(completeTaskVars{
		Path:   "/111",
		Errors: validationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestCompleteTaskHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockCompleteTaskClient
	}{
		{
			Client: &mockCompleteTaskClient{
				GetTaskTypesForDeputyTypeErr: returnedError,
			},
		},
		{
			Client: &mockCompleteTaskClient{
				GetTaskErr: returnedError,
			},
		},
		{
			Client: &mockCompleteTaskClient{
				CompleteTaskErr: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/111", strings.NewReader(""))
			completeTaskReturnedError := renderTemplateForCompleteTask(client, template)(sirius.DeputyDetails{}, w, r)
			assert.Equal(t, returnedError, completeTaskReturnedError)

		})
	}
}
