package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
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
	TaskDetails       model.Task
	Success           bool
	SuccessMessage    string
	Assignees         []model.TeamMember
	IsCurrentAssignee bool
	AppVars
}

func renderTemplateForManageTasks(client ManageTasks, tmpl Template) Handler {
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

		taskDetails.DueDate = sirius.FormatDateTime(sirius.SiriusDate, taskDetails.DueDate, sirius.IsoDate)
		taskDetails.Type = getTaskName(taskDetails.Type, taskTypes)

		assignees, err := client.GetDeputyTeamMembers(ctx, app.DefaultPaTeam, app.DeputyDetails)
		if err != nil {
			return err
		}

		app.PageName = "Manage " + taskDetails.Type + " Task"

		vars := manageTaskVars{
			AppVars:           app,
			TaskDetails:       taskDetails,
			Assignees:         assignees,
			IsCurrentAssignee: true,
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
				updateTaskError := sirius.ValidationErrors{
					"Manage task": {"": "Please update the task information"},
				}

				vars.Errors = util.RenameErrors(updateTaskError)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			err := client.UpdateTask(ctx, app.DeputyId(), taskDetails.Id, dueDate, notes, assigneeId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = RenameErrors(verr.Errors, app.DeputyDetails.DeputyType.Label)
				vars.TaskDetails, vars.IsCurrentAssignee = RetainFormData(vars.TaskDetails, assignees, dueDate, notes, assigneeId)

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/tasks?success=manage&taskType=%s", app.DeputyId(), taskDetails.Type))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func GetAssigneeFromId(id int, teamMembers []model.TeamMember) model.Assignee {
	var teams []model.Team
	var assignee model.TeamMember

	for _, teamMember := range teamMembers {
		if teamMember.ID == id {
			assignee = teamMember
		}
	}

	return model.Assignee{
		Id:          id,
		Teams:       teams,
		DisplayName: assignee.DisplayName,
	}
}

func RetainFormData(task model.Task, assignees []model.TeamMember, dueDate string, notes string, assigneeId int) (model.Task, bool) {
	isCurrentAssignee := true
	task.DueDate = dueDate
	task.Notes = notes

	if task.Assignee.Id != assigneeId {
		task.Assignee = GetAssigneeFromId(assigneeId, assignees)
		isCurrentAssignee = false
	}
	return task, isCurrentAssignee
}

func RenameErrors(errors sirius.ValidationErrors, deputyType string) sirius.ValidationErrors {
	amendedErrors := make(sirius.ValidationErrors)

	for i, s := range errors {
		for k, t := range s {
			if i == "assigneeId" && (k == "notGreaterInclusive" || k == "notLessInclusive") {
				amendedErrors[i] = map[string]string{k: fmt.Sprintf("Enter a name of someone who works on the %s team", deputyType)}
			} else {
				amendedErrors[i] = map[string]string{k: t}
			}
		}
	}
	return amendedErrors
}
