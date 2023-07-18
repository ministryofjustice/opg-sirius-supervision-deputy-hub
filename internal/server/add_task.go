package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type AddTasksClient interface {
	AddTask(ctx sirius.Context, deputyId int, taskType string, dueDate string, notes string, assigneeId int) error
	GetTaskTypes(ctx sirius.Context, deputy sirius.DeputyDetails) ([]model.TaskType, error)
	GetDeputyTeamMembers(ctx sirius.Context, defaultPATeam int, deputy sirius.DeputyDetails) ([]model.TeamMember, error)
	GetTasks(sirius.Context, int) (sirius.TaskList, error)
}

type AddTaskVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	TaskTypes      []model.TaskType
	Assignees      []model.TeamMember
	TaskList       sirius.TaskList
	TaskType       string
	DueDate        string
	Notes          string
	Error          string
	Errors         sirius.ValidationErrors
	SuccessMessage string
}

func renderTemplateForAddTask(client AddTasksClient, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		//routeVars := mux.Vars(r)
		//deputyId := routeVars["id"]
		deputyId := deputyDetails.ID

		taskTypes, err := client.GetTaskTypes(ctx, deputyDetails)
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

			successMessage := ""
			if taskName := r.URL.Query().Get("success"); taskName != "" {
				successMessage = fmt.Sprintf("%s task added", taskName)
			}

			tasklist, err := client.GetTasks(ctx, deputyId)
			if err != nil {
				return err
			}

			vars.TaskList = tasklist
			vars.SuccessMessage = successMessage

			return tmpl.ExecuteTemplate(w, "page", vars)
		} else {
			var (
				taskType   = r.PostFormValue("tasktype")
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

			err := client.AddTask(ctx, deputyId, taskType, dueDate, notes, assigneeId)

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

			return Redirect(fmt.Sprintf("/%d/tasks?success="+taskName, deputyId))
		}
	}
}
