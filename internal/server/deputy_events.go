package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"net/http"
)

type DeputyHubEventInformation interface {
	GetDeputyEvents(sirius.Context, int, int, int) (sirius.TimelineList, error)
}

type deputyHubEventVars struct {
	DeputyEvents sirius.DeputyEvents
	AppVars
	Pagination paginate.Pagination
}

func renderTemplateForDeputyHubEvents(client DeputyHubEventInformation, tmpl Template, vars EnvironmentVars) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		app.PageName = "Timeline"
		params := r.URL.Query()
		page := paginate.GetRequestedPage(params.Get("page"))
		limit := paginate.GetRequestedPage(params.Get("limit"))
		if limit == 1 {
			limit = 25
		}

		perPageOptions := []int{25, 50, 100}
		timelineEventsPerPage := paginate.GetRequestedElementsPerPage(params.Get("limit"), perPageOptions)

		deputyTimelineList, err := client.GetDeputyEvents(ctx, app.DeputyId(), page, limit)
		if err != nil {
			return err
		}

		myUrlBuilder := CreateUrlBuilder(r.URL.Path, timelineEventsPerPage, vars.Prefix)

		pagination := paginate.Pagination{
			CurrentPage:     page,
			TotalPages:      deputyTimelineList.Pages.Total,
			TotalElements:   deputyTimelineList.Total,
			ElementsPerPage: timelineEventsPerPage,
			ElementName:     "timeline event(s)",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      myUrlBuilder,
		}

		vars := deputyHubEventVars{
			DeputyEvents: deputyTimelineList.DeputyEvents,
			Pagination:   pagination,
			AppVars:      app,
		}

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
