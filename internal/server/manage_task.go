package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type ManageTasks interface {
	GetTask(sirius.Context, int) (model.Task, error)
	GetDeputyTeamMembers(ctx sirius.Context, defaultPATeam int, deputy sirius.DeputyDetails) ([]model.TeamMember, error)
	UpdateTask(ctx sirius.Context, deputyId, taskId int, dueDate, notes string, assigneeId int) error
	GetTaskTypesForDeputyType(ctx sirius.Context, deputyType string) ([]model.TaskType, error)
}

type manageTaskVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	TaskDetails    model.Task
	Error          string
	Errors         sirius.ValidationErrors
	Success        bool
	SuccessMessage string
	Assignees      []model.TeamMember
}

func renderTemplateForManageTasks(client ManageTasks, tmpl Template) Handler {
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

		taskDetails.DueDate = sirius.FormatDateTime(sirius.SiriusDate, taskDetails.DueDate, sirius.IsoDate)
		taskDetails.Type = getTaskName(taskDetails.Type, taskTypes)

		defaultPATeam := 1
		assignees, err := client.GetDeputyTeamMembers(ctx, defaultPATeam, deputyDetails)
		if err != nil {
			return err
		}

		vars := manageTaskVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
			TaskDetails:   taskDetails,
			Assignees:     assignees,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var (
				dueDate    = r.PostFormValue("duedate")
				notes      = r.PostFormValue("notes")
				ecm        = r.PostFormValue("assignedto")
				assignedTo = r.PostFormValue("select-assignedto")
			)

			var assigneeId int
			if ecm == "other" {
				assigneeId, _ = strconv.Atoi(assignedTo)
			} else {
				assigneeId, _ = strconv.Atoi(ecm)
			}

			if (dueDate == taskDetails.DueDate) && (notes == taskDetails.Notes) && (assigneeId == taskDetails.Assignee.Id) {
				vars.Errors = sirius.ValidationErrors{
					"Manage task": {"": "Change the page"},
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			err := client.UpdateTask(ctx, deputyDetails.ID, taskDetails.Id, dueDate, notes, assigneeId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := manageTaskVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyDetails: deputyDetails,
					Assignees:     assignees,
					TaskDetails:   taskDetails,
					Errors:        renameErrors(verr.Errors, deputyDetails.DeputyType.Label),
				}

				fmt.Println("error in manage task")
				fmt.Println(verr.Errors)
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/tasks?success=manage"+taskDetails.Type, deputyDetails.ID))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func renameErrors(errors sirius.ValidationErrors, deputyType string) sirius.ValidationErrors {
	amendedErrors := make(sirius.ValidationErrors)

	for i, s := range errors {
		for k, t := range s {
			if i == "assigneeId" && k == "notBetween" {
				amendedErrors[i] = map[string]string{k: fmt.Sprintf("Enter a name of someone who works on the %s team", deputyType)}
			} else {
				amendedErrors[i] = map[string]string{k: t}
			}
		}
	}
	return amendedErrors
}
