package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type TasksClient interface {
}

type TasksVars struct {
	DeputyDetails  sirius.DeputyDetails
	Error          string
	SuccessMessage string
}

func renderTemplateForTasksTab(client TasksClient, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodGet:
			successMessage := ""
			if r.URL.Query().Get("success") == "true" {
				successMessage = "Task added"
			}

			vars := deputyHubNotesVars{
				Path:           r.URL.Path,
				DeputyDetails:  deputyDetails,
				SuccessMessage: successMessage,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
