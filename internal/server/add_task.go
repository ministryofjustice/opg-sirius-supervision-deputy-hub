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
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	TaskTypes     []model.TaskType
	Assignees     []model.TeamMember
	TaskType      string
	DueDate       string
	Notes         string
	Error         string
	Errors        sirius.ValidationErrors
	IsManageTasks bool
}

func renderTemplateForAddTask(client AddTasksClient, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)

		taskTypes, err := client.GetTaskTypesForDeputyType(ctx, deputyDetails.DeputyType.Handle)
		if err != nil {
			return err
		}

		defaultPATeam := 1
		assignees, err := client.GetDeputyTeamMembers(ctx, defaultPATeam, deputyDetails)
		if err != nil {
			return err
		}

		vars := AddTaskVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			TaskTypes:     taskTypes,
			DeputyDetails: deputyDetails,
			Assignees:     assignees,
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

			err := client.AddTask(ctx, deputyDetails.ID, taskType, typeName, dueDate, notes, assigneeId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars = AddTaskVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					TaskTypes:     taskTypes,
					DeputyDetails: deputyDetails,
					Assignees:     assignees,
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

			return Redirect(fmt.Sprintf("/%d/tasks?success="+taskName, deputyDetails.ID))
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
