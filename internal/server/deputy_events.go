package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
)

type DeputyHubEventInformation interface {
	GetDeputyEvents(sirius.Context, int, int, int) (sirius.TimelineList, error)
}

type deputyHubEventVars struct {
	DeputyEvents sirius.DeputyEvents
	ListPage
}

func renderTemplateForDeputyHubEvents(client DeputyHubEventInformation, tmpl Template, vars EnvironmentVars) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		vars := deputyHubEventVars{}
		app.PageName = "Timeline"
		vars.AppVars = app

		params := r.URL.Query()
		page := paginate.GetRequestedPage(params.Get("page"))
		perPageOptions := []int{25, 50, 100}
		perPage := paginate.GetRequestedElementsPerPage(params.Get("limit"), perPageOptions)

		vars.PerPage = perPage

		deputyTimelineList, err := client.GetDeputyEvents(ctx, app.DeputyId(), page, perPage)
		if err != nil {
			return err
		}

		myUrlBuilder := CreateUrlBuilder(r.URL.Path, perPage, vars.Prefix)

		vars.Pagination = paginate.Pagination{
			CurrentPage:     page,
			TotalPages:      deputyTimelineList.Pages.Total,
			TotalElements:   deputyTimelineList.Total,
			ElementsPerPage: vars.PerPage,
			ElementName:     "timeline event(s)",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      myUrlBuilder,
		}

		vars.DeputyEvents = deputyTimelineList.DeputyEvents

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func CreateUrlBuilder(urlPath string, timelineEventsPerPage int, prefix string) urlbuilder.UrlBuilder {
	path := prefix + urlPath
	return urlbuilder.UrlBuilder{
		OriginalPath:    path,
		SelectedPerPage: timelineEventsPerPage,
	}
}
