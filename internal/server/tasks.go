package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type TasksClient interface {
	GetTaskTypesForDeputyType(ctx sirius.Context, deputyType string) ([]model.TaskType, error)
	GetTasks(ctx sirius.Context, deputyId int) (sirius.TaskList, error)
}

type TasksVars struct {
	TaskTypes      []model.TaskType
	TaskList       sirius.TaskList
	TaskType       string
	DueDate        string
	Notes          string
	SuccessMessage string
	AppVars
}

func renderTemplateForTasks(client TasksClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)

		taskTypes, err := client.GetTaskTypesForDeputyType(ctx, app.DeputyDetails.DeputyType.Handle)
		if err != nil {
			return err
		}

		taskType := r.URL.Query().Get("taskType")

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

		taskList, err := client.GetTasks(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		vars := TasksVars{
			AppVars:        app,
			TaskTypes:      taskTypes,
			TaskList:       taskList,
			SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
