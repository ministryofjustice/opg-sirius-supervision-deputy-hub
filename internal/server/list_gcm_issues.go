package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
)

type GetGcmIssues interface {
	GetGCMIssues(ctx sirius.Context, deputyId int, params sirius.GcmIssuesParams) (sirius.GcmIssueList, error)
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
	PerPage        int
	Pagination     paginate.Pagination
}

func (gv GcmIssuesVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		OriginalPath:    fmt.Sprintf("%s-issues", gv.GCMIssueStatus),
		SelectedPerPage: gv.PerPage,
		SelectedSort:    gv.Sort,
	}
}

func renderTemplateForGcmIssues(client GetGcmIssues, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		var successMessage string
		switch r.Method {
		case http.MethodGet:
			urlParams := r.URL.Query()
			page := paginate.GetRequestedPage(urlParams.Get("page"))
			perPageOptions := []int{25, 50, 100}
			perPage := paginate.GetRequestedElementsPerPage(urlParams.Get("limit"), perPageOptions)
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
				Limit:       perPage,
				Page:        page,
			}

			gcmIssueList, err := client.GetGCMIssues(ctx, app.DeputyId(), params)
			if err != nil {
				return err
			}

			vars := GcmIssuesVars{
				SuccessMessage: successMessage,
				AppVars:        app,
				GcmIssues:      gcmIssueList.GcmIssues,
				Sort:           sort,
				GCMIssueStatus: issueStatus,
				PerPage:        perPage,
			}

			vars.UrlBuilder = vars.CreateUrlBuilder()

			if page > gcmIssueList.Pages.PageTotal && gcmIssueList.Pages.PageTotal > 0 {
				return Redirect(vars.UrlBuilder.GetPaginationUrl(gcmIssueList.Pages.PageTotal, perPage))
			}

			vars.Pagination = paginate.Pagination{
				CurrentPage:     gcmIssueList.Pages.PageCurrent,
				TotalPages:      gcmIssueList.Pages.PageTotal,
				TotalElements:   gcmIssueList.Total,
				ElementsPerPage: perPage,
				ElementName:     "issues",
				PerPageOptions:  perPageOptions,
				UrlBuilder:      vars.UrlBuilder,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:

			// Specify max file size to 10mb
			r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

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
