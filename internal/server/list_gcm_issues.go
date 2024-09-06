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
	CloseGCMIssues(ctx sirius.Context, gcmIds []string) error
}

type GcmIssuesVars struct {
	ErrorMessage   string
	SuccessMessage string
	GcmIssues      []sirius.GcmIssue
	Sort           urlbuilder.Sort
	AppVars
	UrlBuilder     urlbuilder.UrlBuilder
	GCMIssueStatus string
}

func (gv GcmIssuesVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		OriginalPath: "open-issues",
		SelectedSort: gv.Sort,
	}
}

func renderTemplateForGcmIssues(client GetGcmIssues, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		var successMessage string
		switch r.Method {
		case http.MethodGet:
			urlParams := r.URL.Query()
			sort := urlbuilder.CreateSortFromURL(urlParams, []string{"createdDate", "issueType"})
			path := r.URL.Path
			issueStatus := ""

			if strings.Contains(path, "open-issues") {
				issueStatus = "open"
			} else if strings.Contains(path, "closed-issues") {
				issueStatus = "closed"
			}

			switch r.URL.Query().Get("success") {
			case "addGcmIssue":
				successMessage = "GCM Issue added"
			case "closedGcms":
				selectedGCMCount := r.URL.Query().Get("count")
				successMessage = fmt.Sprintf("You have closed %s number(s) of GCM issues.", selectedGCMCount)
			default:
				successMessage = ""
			}

			app.PageName = "General Case Manager issues"
			params := sirius.GcmIssuesParams{
				IssueStatus: fmt.Sprintf("%s:%s", "status", issueStatus),
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
				GCMIssueStatus: issueStatus,
			}

			vars.UrlBuilder = vars.CreateUrlBuilder()

			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				return err
			}
			selectedGCMs := r.Form["selected-gcms"]

			err = client.CloseGCMIssues(ctx, selectedGCMs)
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/gcm-issues/open-issues?success=closedGcms&count=%d", app.DeputyId(), len(selectedGCMs)))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
