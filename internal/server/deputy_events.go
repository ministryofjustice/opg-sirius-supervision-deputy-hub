package server

import (
"github.com/gorilla/mux"
"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
"net/http"
"strconv"
)

type DeputyHubEventInformation interface {
	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
	GetDeputyEvents(sirius.Context, int) (sirius.DeputyEvents, error)
}

type deputyHubEventVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	DeputyEvents sirius.DeputyEvents
	Error         string
	Errors        sirius.ValidationErrors
}

func renderTemplateForDeputyHubEvents(client DeputyHubEventInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails, err := client.GetDeputyDetails(ctx, deputyId)
		deputyEvents, err := client.GetDeputyEvents(ctx, deputyId)
		if err != nil {
			return err
		}

		vars := deputyHubEventVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
			DeputyEvents: deputyEvents,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
