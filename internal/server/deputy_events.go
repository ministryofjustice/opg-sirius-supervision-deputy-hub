package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type DeputyHubEventInformation interface {
	GetDeputyEvents(sirius.Context, sirius.DeputyDetails) (sirius.DeputyEvents, error)
}

type deputyHubEventVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	DeputyEvents  sirius.DeputyEvents
	Error         string
}

func renderTemplateForDeputyHubEvents(client DeputyHubEventInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		deputyEvents, err := client.GetDeputyEvents(ctx, deputyDetails)
		if err != nil {
			return err
		}

		vars := deputyHubEventVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
			DeputyEvents:  deputyEvents,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
