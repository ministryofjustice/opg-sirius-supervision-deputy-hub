package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubEventInformation interface {
	GetDeputyEvents(sirius.Context, int) (sirius.DeputyEventCollection, error)
}

type deputyHubEventVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	DeputyEvents  sirius.DeputyEventCollection
	Error         string
}

func renderTemplateForDeputyHubEvents(client DeputyHubEventInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyEvents, err := client.GetDeputyEvents(ctx, deputyId)
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
