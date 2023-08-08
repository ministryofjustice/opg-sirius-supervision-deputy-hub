package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type DeputyHubEventInformation interface {
	GetDeputyEvents(sirius.Context, int) (sirius.DeputyEvents, error)
}

type deputyHubEventVars struct {
	DeputyEvents sirius.DeputyEvents
	AppVars
}

func renderTemplateForDeputyHubEvents(client DeputyHubEventInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		deputyEvents, err := client.GetDeputyEvents(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		vars := deputyHubEventVars{
			DeputyEvents: deputyEvents,
			AppVars:      app,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
