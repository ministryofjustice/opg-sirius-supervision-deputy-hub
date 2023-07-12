package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type TasksClient interface {
	GetTaskTypes(ctx sirius.Context, deputy sirius.DeputyDetails) ([]model.TaskType, error)
	GetTasks(sirius.Context, string) (sirius.TaskList, error)
}

type tasksVars struct {
	DeputyDetails  sirius.DeputyDetails
	Error          string
	SuccessMessage string
	Path           string
	TaskTypes      []model.TaskType
	Tasklist       sirius.TaskList
}

func renderTemplateForTasksTab(client TasksClient, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodGet:
			successMessage := ""
			if taskName := r.URL.Query().Get("success"); taskName != "" {
				successMessage = fmt.Sprintf("%s task added", taskName)
			}
			ctx := getContext(r)
			routeVars := mux.Vars(r)
			deputyId := routeVars["id"]

			taskTypes, err := client.GetTaskTypes(ctx, deputyDetails)
			if err != nil {
				return err
			}

			tasklist, err := client.GetTasks(ctx, deputyId)
			if err != nil {
				return err
			}

			vars := tasksVars{
				Path:           r.URL.Path,
				DeputyDetails:  deputyDetails,
				SuccessMessage: successMessage,
				TaskTypes:      taskTypes,
				Tasklist:       tasklist,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
