package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type TasksClient interface {
	GetTaskTypesForDeputyType(ctx sirius.Context, deputyType string) ([]model.TaskType, error)
	GetDeputyTeamMembers(ctx sirius.Context, defaultPATeam int, deputy sirius.DeputyDetails) ([]model.TeamMember, error)
	GetTasks(sirius.Context, int) (sirius.TaskList, error)
}

type TasksVars struct {
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

func renderTemplateForTasks(client TasksClient, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)

		deputyId := deputyDetails.ID

		taskTypes, err := client.GetTaskTypesForDeputyType(ctx, deputyDetails.DeputyType.Handle)
		if err != nil {
			return err
		}

		defaultPATeam := 1
		assignees, err := client.GetDeputyTeamMembers(ctx, defaultPATeam, deputyDetails)
		if err != nil {
			return err
		}

		successMessage := ""
		if taskName := r.URL.Query().Get("success"); taskName != "" {
			successMessage = fmt.Sprintf("%s task added", taskName)
		}

		taskList, err := client.GetTasks(ctx, deputyId)
		if err != nil {
			return err
		}

		vars := TasksVars{
			Path:           r.URL.Path,
			XSRFToken:      ctx.XSRFToken,
			TaskTypes:      taskTypes,
			DeputyDetails:  deputyDetails,
			Assignees:      assignees,
			TaskList:       taskList,
			SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
