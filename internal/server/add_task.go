package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type AddTasksClient interface {
	AddTask(ctx sirius.Context, deputyId int, taskType string, typeName string, dueDate string, notes string, assigneeId int) error
	GetTaskTypesForDeputyType(ctx sirius.Context, deputyType string) ([]model.TaskType, error)
	GetDeputyTeamMembers(ctx sirius.Context, defaultPATeam int, deputy sirius.DeputyDetails) ([]model.TeamMember, error)
}

type AddTaskVars struct {
	TaskTypes      []model.TaskType
	Assignees      []model.TeamMember
	TaskType       string
	DueDate        string
	Notes          string
	SuccessMessage string
	IsManageTasks  bool
	AppVars
}

func renderTemplateForAddTask(client AddTasksClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)

		taskTypes, err := client.GetTaskTypesForDeputyType(ctx, app.DeputyDetails.DeputyType.Handle)
		if err != nil {
			return err
		}

		assignees, err := client.GetDeputyTeamMembers(ctx, app.DefaultPaTeam, app.DeputyDetails)
		if err != nil {
			return err
		}

		vars := AddTaskVars{
			TaskTypes: taskTypes,
			Assignees: assignees,
			AppVars:   app,
		}

		if r.Method == http.MethodGet {
			return tmpl.ExecuteTemplate(w, "page", vars)
		} else {
			var (
				taskType   = r.PostFormValue("tasktype")
				typeName   = getTaskName(taskType, taskTypes)
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

			err := client.AddTask(ctx, app.DeputyId(), taskType, typeName, dueDate, notes, assigneeId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.TaskTypes = taskTypes
				vars.Assignees = assignees
				vars.TaskType = taskType
				vars.DueDate = dueDate
				vars.Notes = notes
				vars.Errors = verr.Errors
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

			return Redirect(fmt.Sprintf("/%d/tasks?success=add&taskType=%s", app.DeputyId(), taskName))
		}
	}
}

func getTaskName(handle string, types []model.TaskType) string {
	for _, t := range types {
		if handle == t.Handle {
			return t.Description
		}
	}
	return ""
}
