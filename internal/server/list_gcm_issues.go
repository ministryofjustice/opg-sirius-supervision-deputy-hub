package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type GetGcmIssues interface {
	GetGCMIssues(ctx sirius.Context, deputyId int) ([]sirius.GcmIssue, error)
}

type GcmIssuesVars struct {
	ErrorMessage   string
	SuccessMessage string
	GcmIssues      []sirius.GcmIssue
	AppVars
}

func renderTemplateForGcmIssues(client GetGcmIssues, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addGcmIssue":
			successMessage = "GCM Issue added"
		default:
			successMessage = ""
		}

		app.PageName = "General Case Manager issues"
		gcmIssues, err := client.GetGCMIssues(ctx, app.DeputyId())

		if err != nil {
			return err
		}

		vars := GcmIssuesVars{
			SuccessMessage: successMessage,
			AppVars:        app,
			GcmIssues:      gcmIssues,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
