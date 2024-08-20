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
		//urlParams := r.URL.Query()
		//sort := urlbuilder.CreateSortFromURL(urlParams, []string{"createdDate", "issueType"})

		ctx := getContext(r)
		//path := r.URL.Path
		//issueStatus := ""
		//
		//if strings.Contains(path, "open-issues") {
		//	issueStatus = "open"
		//} else if strings.Contains(path, "resolved-issues") {
		//	issueStatus = "resolved"
		//}

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addGcmIssue":
			successMessage = "GCM Issue added"
		default:
			successMessage = ""
		}

		app.PageName = "General Case Manager issues"
		//params := sirius.GcmIssuesParams{
		//	IssueStatus: issueStatus,
		//	//Sort:        fmt.Sprintf("%s:%s", sort.OrderBy, sort.GetDirection()),
		//}

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
