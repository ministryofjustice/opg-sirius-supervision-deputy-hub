package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strconv"
)

type CompleteTask interface {
	GetTask(sirius.Context, int) (model.Task, error)
	GetTaskTypesForDeputyType(ctx sirius.Context, deputyType string) ([]model.TaskType, error)
	CompleteTask(sirius.Context, int, string) error
}

type completeTaskVars struct {
	TaskDetails    model.Task
	CompletedNotes string
	SuccessMessage string
	AppVars
}

type CompleteTaskHandler struct {
	router
}

func (h *CompleteTaskHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	taskId, _ := strconv.Atoi(r.PathValue("taskId"))

	taskTypes, err := h.Client().GetTaskTypesForDeputyType(ctx, v.DeputyType())
	if err != nil {
		return err
	}

	taskDetails, err := h.Client().GetTask(ctx, taskId)
	if err != nil {
		return err
	}

	taskDetails.Type = getTaskName(taskDetails.Type, taskTypes)

	v.PageName = "Complete Task"

	vars := completeTaskVars{
		AppVars:     v,
		TaskDetails: taskDetails,
	}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars, vars.AppVars)

	case http.MethodPost:
		var (
			notes = r.PostFormValue("notes")
		)

		err = h.Client().CompleteTask(ctx, taskDetails.Id, notes)

		if verr, ok := err.(sirius.ValidationError); ok {

			vars.Errors = util.RenameErrors(verr.Errors)
			vars.CompletedNotes = notes

			w.WriteHeader(http.StatusBadRequest)
			return h.execute(w, r, vars, vars.AppVars)
		}
		if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d/tasks?success=complete&taskType=%s", v.DeputyId(), taskDetails.Type))

	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
