package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strings"
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
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)

		taskTypes, err := client.GetTaskTypesForDeputyType(ctx, appVars.DeputyDetails.DeputyType.Handle)
		if err != nil {
			return err
		}

		successMessage := ""

		if taskName := r.URL.Query().Get("success"); taskName != "" {
			if strings.Contains(taskName, "manage") {
				successMessage = fmt.Sprintf("%s task updated", strings.ReplaceAll(taskName, "manage", ""))
			} else {
				successMessage = fmt.Sprintf("%s task added", taskName)
			}
		}

		taskList, err := client.GetTasks(ctx, appVars.DeputyId())
		if err != nil {
			return err
		}

		vars := TasksVars{
			AppVars:        appVars,
			TaskTypes:      taskTypes,
			TaskList:       taskList,
			SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
