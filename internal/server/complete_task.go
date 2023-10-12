package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
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

func renderTemplateForCompleteTask(client CompleteTask, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		taskId, _ := strconv.Atoi(routeVars["taskId"])

		taskTypes, err := client.GetTaskTypesForDeputyType(ctx, app.DeputyType())
		if err != nil {
			return err
		}

		taskDetails, err := client.GetTask(ctx, taskId)
		if err != nil {
			return err
		}

		taskDetails.Type = getTaskName(taskDetails.Type, taskTypes)

		app.PageName = "Complete Task"

		vars := completeTaskVars{
			AppVars:     app,
			TaskDetails: taskDetails,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var (
				notes = r.PostFormValue("notes")
			)

			err = client.CompleteTask(ctx, taskDetails.Id, notes)

			if verr, ok := err.(sirius.ValidationError); ok {

				vars.Errors = verr.Errors
				vars.CompletedNotes = notes

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/tasks?success=complete&taskType=%s", app.DeputyId(), taskDetails.Type))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
