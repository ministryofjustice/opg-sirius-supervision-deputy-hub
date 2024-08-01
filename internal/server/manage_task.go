package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strconv"
)

type manageTaskVars struct {
	TaskDetails       model.Task
	Success           bool
	SuccessMessage    string
	Assignees         []model.TeamMember
	IsCurrentAssignee bool
	AppVars
}

type ManageTaskHandler struct {
	router
}

func (h *ManageTaskHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
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

	taskDetails.DueDate = sirius.FormatDateTime(sirius.SiriusDate, taskDetails.DueDate, sirius.IsoDate)
	taskDetails.Type = getTaskName(taskDetails.Type, taskTypes)

	assignees, err := h.Client().GetDeputyTeamMembers(ctx, v.DefaultPaTeam, v.DeputyDetails)
	if err != nil {
		return err
	}

	v.PageName = "Manage " + taskDetails.Type + " Task"

	vars := manageTaskVars{
		AppVars:           v,
		TaskDetails:       taskDetails,
		Assignees:         assignees,
		IsCurrentAssignee: true,
	}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars)

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
			return h.execute(w, r, vars)
		}

		err := h.Client().UpdateTask(ctx, v.DeputyId(), taskDetails.Id, dueDate, notes, assigneeId)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = RenameErrors(verr.Errors, v.DeputyDetails.DeputyType.Label)
			vars.TaskDetails, vars.IsCurrentAssignee = RetainFormData(vars.TaskDetails, assignees, dueDate, notes, assigneeId)

			w.WriteHeader(http.StatusBadRequest)
			return h.execute(w, r, vars)
		}
		if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d/tasks?success=manage&taskType=%s", v.DeputyId(), taskDetails.Type))

	default:
		return StatusError(http.StatusMethodNotAllowed)
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
			if i == "assigneeId" && k == "notBetween" {
				amendedErrors[i] = map[string]string{k: fmt.Sprintf("Enter a name of someone who works on the %s team", deputyType)}
			} else {
				amendedErrors[i] = map[string]string{k: t}
			}
		}
	}
	return amendedErrors
}
