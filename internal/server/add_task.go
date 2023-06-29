package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type AddTasksClient interface {
	AddTask(ctx sirius.Context, deputyId int, taskType string, dueDate string, notes string) error
	GetTaskTypes(ctx sirius.Context, deputy sirius.DeputyDetails) ([]sirius.TaskType, error)
}

type AddTaskVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	TaskTypes     []sirius.TaskType
	TaskType      string
	DueDate       string
	Notes         string
	Error         string
	Errors        sirius.ValidationErrors
}

func renderTemplateForAddTask(client AddTasksClient, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		taskTypes, err := client.GetTaskTypes(ctx, deputyDetails)
		if err != nil {
			return err
		}

		vars := AddTaskVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			TaskTypes:     taskTypes,
			DeputyDetails: deputyDetails,
		}

		if r.Method == http.MethodGet {
			return tmpl.ExecuteTemplate(w, "page", vars)
		} else {
			var (
				taskType = r.PostFormValue("tasktype")
				dueDate  = r.PostFormValue("duedate")
				notes    = r.PostFormValue("notes")
			)

			err := client.AddTask(ctx, deputyId, taskType, dueDate, notes)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars = AddTaskVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					TaskTypes:     taskTypes,
					DeputyDetails: deputyDetails,
					TaskType:      taskType,
					DueDate:       dueDate,
					Notes:         notes,
					Errors:        verr.Errors,
				}
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			var taskName string
			for _, t := range taskTypes {
				if t.Handle == taskType {
					taskName = t.Description
				}
			}

			return Redirect(fmt.Sprintf("/%d/tasks?success="+taskName, deputyId))
		}
	}
}
