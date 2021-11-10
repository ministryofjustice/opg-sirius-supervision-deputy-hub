package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubEventInformation interface {
	GetDeputyDetails(sirius.Context, string, int) (sirius.DeputyDetails, error)
	GetDeputyEvents(sirius.Context, int) (sirius.DeputyEventCollection, error)
}

type deputyHubEventVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	DeputyEvents  sirius.DeputyEventCollection
	Error         string
	ErrorMessage  string
	Errors        sirius.ValidationErrors
}

func renderTemplateForDeputyHubEvents(client DeputyHubEventInformation, defaultPATeam string, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails, err := client.GetDeputyDetails(ctx, defaultPATeam, deputyId)
		if err != nil {
			return err
		}
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

		if vars.DeputyDetails.ExecutiveCaseManager.EcmName == defaultPATeam {
			vars.ErrorMessage = "An executive case manager has not been assigned. "
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
