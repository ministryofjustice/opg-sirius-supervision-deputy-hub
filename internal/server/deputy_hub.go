package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type DeputyHubInformation interface {
	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
}

type deputyHubVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	ErrorMessage  string
	Errors        sirius.ValidationErrors
}

func renderTemplateForDeputyHub(client DeputyHubInformation, defaultPATeam string, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails, err := client.GetDeputyDetails(ctx, deputyId)
		if err != nil {
			return err
		}

		vars := deputyHubVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
		}

        if vars.DeputyDetails.OrganisationTeamOrDepartmentName == defaultPATeam {
            vars.ErrorMessage = "An executive case manager has not been assigned. "
        }

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
