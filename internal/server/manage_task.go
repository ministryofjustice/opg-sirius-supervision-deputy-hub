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
	EditTask(ctx sirius.Context, deputyId, taskId int, dueDate, notes string, assigneeId int) error
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
		//deputyId, _ := strconv.Atoi(routeVars["id"])
		taskId, _ := strconv.Atoi(routeVars["taskId"])

		taskDetails, err := client.GetTask(ctx, taskId)
		if err != nil {
			return err
		}

		taskDetails.DueDate = sirius.FormatDateTime(sirius.SiriusDate, taskDetails.DueDate, sirius.IsoDate)

		defaultPATeam := 1
		assignees, err := client.GetDeputyTeamMembers(ctx, defaultPATeam, deputyDetails)
		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:

			vars := manageTaskVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
				TaskDetails:   taskDetails,
				Assignees:     assignees,
			}

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

			err := client.EditTask(ctx, deputyDetails.ID, taskDetails.Id, dueDate, notes, assigneeId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := manageTaskVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyDetails: deputyDetails,
					Assignees:     assignees,
					TaskDetails:   taskDetails,
					Errors:        verr.Errors,
				}
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/tasks", deputyDetails.ID))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
