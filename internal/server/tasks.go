package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type TasksVars struct {
	TaskTypes      []model.TaskType
	TaskList       sirius.TaskList
	TaskType       string
	DueDate        string
	Notes          string
	SuccessMessage string
	AppVars
}

type TasksHandler struct {
	router
}

func (h *TasksHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	taskTypes, err := h.Client().GetTaskTypesForDeputyType(ctx, v.DeputyType())
	if err != nil {
		return err
	}

	taskType := r.URL.Query().Get("taskType")

	taskList, err := h.Client().GetTasks(ctx, v.DeputyId())
	if err != nil {
		return err
	}

	var successMessage string

	switch r.URL.Query().Get("success") {
	case "add":
		successMessage = fmt.Sprintf("%s task added", taskType)
	case "manage":
		successMessage = fmt.Sprintf("%s task updated", taskType)
	case "complete":
		successMessage = fmt.Sprintf("%s task completed", taskType)
	default:
		successMessage = ""
	}

	v.PageName = "Deputy tasks"

	vars := TasksVars{
		AppVars:        v,
		TaskTypes:      taskTypes,
		TaskList:       taskList,
		SuccessMessage: successMessage,
	}

	return h.execute(w, r, vars)
}
