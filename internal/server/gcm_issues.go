package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type GetGcmIssuesClient interface {
	GetAssurances(ctx sirius.Context, deputyId int) ([]model.Assurance, error)
}

type GcmIssuesVars struct {
	ErrorMessage   string
	SuccessMessage string
	AppVars
}

func renderTemplateForGcmIssues(client GetAssurancesClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		//ctx := getContext(r)

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addGcmIssue":
			successMessage = "GCM Issue added"
		default:
			successMessage = ""
		}

		app.PageName = "General Case Manager issues"

		vars := GcmIssuesVars{
			SuccessMessage: successMessage,
			AppVars:        app,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
