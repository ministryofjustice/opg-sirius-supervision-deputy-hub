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
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
	CompleteTask(sirius.Context, int, int, string) error
}

type completeTaskVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	TaskDetails    model.Task
	Error          string
	Errors         sirius.ValidationErrors
	Success        bool
	SuccessMessage string
}

func renderTemplateForCompleteTask(client CompleteTask, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		taskId, _ := strconv.Atoi(routeVars["taskId"])

		taskTypes, err := client.GetTaskTypesForDeputyType(ctx, deputyDetails.DeputyType.Handle)
		if err != nil {
			return err
		}

		taskDetails, err := client.GetTask(ctx, taskId)
		if err != nil {
			return err
		}

		taskDetails.Type = getTaskName(taskDetails.Type, taskTypes)

		vars := completeTaskVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
			TaskDetails:   taskDetails,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var (
				notes = r.PostFormValue("notes")
			)

			fmt.Println("notes")
			fmt.Println(notes)

			userDetails, err := client.GetUserDetails(ctx)
			if err != nil {
				return err
			}
			fmt.Println("after user details")
			err = client.CompleteTask(ctx, userDetails.ID, taskDetails.Id, notes)
			fmt.Println("after complete details")

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := completeTaskVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyDetails: deputyDetails,
					TaskDetails:   taskDetails,
					Errors:        verr.Errors,
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/tasks?success=complete&type="+taskDetails.Type, deputyDetails.ID))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
