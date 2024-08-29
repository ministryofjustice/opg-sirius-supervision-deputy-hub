package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"net/http"
	"strings"
)

type GetGcmIssues interface {
	GetGCMIssues(ctx sirius.Context, deputyId int, params sirius.GcmIssuesParams) ([]sirius.GcmIssue, error)
}

type GcmIssuesVars struct {
	ErrorMessage   string
	SuccessMessage string
	GcmIssues      []sirius.GcmIssue
	Sort           urlbuilder.Sort
	AppVars
	UrlBuilder urlbuilder.UrlBuilder
}

func (gv GcmIssuesVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		OriginalPath: "open-issues",
		SelectedSort: gv.Sort,
	}
}

func renderTemplateForGcmIssues(client GetGcmIssues, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		urlParams := r.URL.Query()
		sort := urlbuilder.CreateSortFromURL(urlParams, []string{"createdDate", "issueType"})

		ctx := getContext(r)
		path := r.URL.Path
		issueStatus := ""

		if strings.Contains(path, "open-issues") {
			issueStatus = "open"
		} else if strings.Contains(path, "resolved-issues") {
			issueStatus = "resolved"
		}

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addGcmIssue":
			successMessage = "GCM Issue added"
		default:
			successMessage = ""
		}

		app.PageName = "General Case Manager issues"
		params := sirius.GcmIssuesParams{
			IssueStatus: issueStatus,
			Sort:        fmt.Sprintf("%s:%s", sort.OrderBy, sort.GetDirection()),
		}

		gcmIssues, err := client.GetGCMIssues(ctx, app.DeputyId(), params)

		if err != nil {
			return err
		}

		vars := GcmIssuesVars{
			SuccessMessage: successMessage,
			AppVars:        app,
			GcmIssues:      gcmIssues,
			Sort:           sort,
		}

		vars.UrlBuilder = vars.CreateUrlBuilder()

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
