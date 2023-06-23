package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddTasksClient struct {
	count       int
	lastCtx     sirius.Context
	err         error
	userDetails sirius.UserDetails
	taskTypes   []sirius.TaskType
}

func (m *mockAddTasksClient) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetails, m.err
}

func (m *mockAddTasksClient) AddTask(ctx sirius.Context, deputyId int, taskType string, dueDate string, notes string) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockAddTasksClient) GetTaskTypes(ctx sirius.Context, details sirius.DeputyDetails) ([]sirius.TaskType, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskTypes, m.err
}

func TestLoadAddTaskForm(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTasksClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddTask(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestAddTask_success(t *testing.T) {
	assert := assert.New(t)
	client := &mockAddTasksClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddTask(client, nil)(sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/123?success=true"))
}

//func TestAddFirmValidationErrors(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockFirmInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"stringLengthTooLong": "The firm name must be 255 characters or fewer",
//		},
//	}
//
//	client.err = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
//	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForAddFirm(client, template)(sirius.DeputyDetails{}, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//
//	assert.Equal(addFirmVars{
//		Path:   "/133",
//		Errors: validationErrors,
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}
//
//func TestErrorAddFirmMessageWhenIsEmpty(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockFirmInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"isEmpty": "The firm name is required and can't be empty",
//		},
//	}
//
//	client.err = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
//	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForAddFirm(client, template)(sirius.DeputyDetails{}, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//
//	expectedValidationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"isEmpty": "The firm name is required and can't be empty",
//		},
//	}
//
//	assert.Equal(addFirmVars{
//		Path:   "/133",
//		Errors: expectedValidationErrors,
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}
