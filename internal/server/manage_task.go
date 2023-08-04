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
	TaskDetails    model.Task
	Success        bool
	SuccessMessage string
	Assignees      []model.TeamMember
	AppVars
}

func renderTemplateForManageTasks(client ManageTasks, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		taskId, _ := strconv.Atoi(routeVars["taskId"])

		taskTypes, err := client.GetTaskTypesForDeputyType(ctx, appVars.DeputyDetails.DeputyType.Handle)
		if err != nil {
			return err
		}

		taskDetails, err := client.GetTask(ctx, taskId)
		if err != nil {
			return err
		}

		taskDetails.DueDate = sirius.FormatDateTime(sirius.SiriusDate, taskDetails.DueDate, sirius.IsoDate)
		taskDetails.Type = getTaskName(taskDetails.Type, taskTypes)

		assignees, err := client.GetDeputyTeamMembers(ctx, appVars.DefaultPaTeam, appVars.DeputyDetails)
		if err != nil {
			return err
		}

		vars := manageTaskVars{
			AppVars:     appVars,
			TaskDetails: taskDetails,
			Assignees:   assignees,
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
					"Manage task": {"": "Please update the task information"},
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			err := client.UpdateTask(ctx, appVars.DeputyId(), taskDetails.Id, dueDate, notes, assigneeId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = RenameErrors(verr.Errors, appVars.DeputyDetails.DeputyType.Label)

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/tasks?success=manage&taskType=%s", appVars.DeputyId(), taskDetails.Type))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func RenameErrors(errors sirius.ValidationErrors, deputyType string) sirius.ValidationErrors {
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
