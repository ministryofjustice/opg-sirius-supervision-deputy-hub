package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockTasksClient struct {
	count           int
	lastCtx         sirius.Context
	GetTaskTypesErr error
	GetTasksErr     error
	taskTypes       []model.TaskType
	teamMembers     []model.TeamMember
	tasks           sirius.TaskList
}

func (m *mockTasksClient) GetTaskTypesForDeputyType(ctx sirius.Context, deputyId string) ([]model.TaskType, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskTypes, m.GetTaskTypesErr
}

func (m *mockTasksClient) GetTasks(ctx sirius.Context, deputyId int) (sirius.TaskList, error) {
	m.count += 1

	return m.tasks, m.GetTasksErr
}

func TestNavigateTasksTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockTasksClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForTasks(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestTasksHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockTasksClient
	}{
		{
			Client: &mockTasksClient{
				GetTaskTypesErr: returnedError,
			},
		},
		{
			Client: &mockTasksClient{
				GetTasksErr: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/123", strings.NewReader(""))

			addFirmReturnedError := renderTemplateForTasks(client, template)(sirius.DeputyDetails{}, w, r)
			assert.Equal(t, returnedError, addFirmReturnedError)
		})
	}
}
