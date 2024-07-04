package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strconv"
)

type AddTask interface {
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

type AddTaskHandler struct {
	router
}

func (h *AddTaskHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	taskTypes, err := h.Client().GetTaskTypesForDeputyType(ctx, v.DeputyType())
	if err != nil {
		return err
	}

	assignees, err := h.Client().GetDeputyTeamMembers(ctx, v.DefaultPaTeam, v.DeputyDetails)
	if err != nil {
		return err
	}

	v.PageName = "Add a deputy task"

	vars := AddTaskVars{
		TaskTypes: taskTypes,
		Assignees: assignees,
		AppVars:   v,
	}

	if r.Method == http.MethodGet {
		return h.execute(w, r, vars, v)
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

		err := h.Client().AddTask(ctx, v.DeputyId(), taskType, typeName, dueDate, notes, assigneeId)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.TaskTypes = taskTypes
			vars.Assignees = assignees
			vars.TaskType = taskType
			vars.DueDate = dueDate
			vars.Notes = notes
			vars.Errors = util.RenameErrors(verr.Errors)
			w.WriteHeader(http.StatusBadRequest)
			return h.execute(w, r, vars, v)
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

		return Redirect(fmt.Sprintf("/%d/tasks?success=add&taskType=%s", v.DeputyId(), taskName))
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
