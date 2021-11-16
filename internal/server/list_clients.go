package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubClientInformation interface {
	GetDeputyClients(sirius.Context, int, string, string) (sirius.DeputyClientDetails, sirius.AriaSorting, error)
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
}

type listClientsVars struct {
	Path                 string
	XSRFToken            string
	AriaSorting          sirius.AriaSorting
	DeputyClientsDetails sirius.DeputyClientDetails
	DeputyDetails        sirius.DeputyDetails
	Error                string
	ErrorMessage         string
	Errors               sirius.ValidationErrors
}

func renderTemplateForClientTab(client DeputyHubClientInformation, defaultPATeam int, tmpl Template) Handler {
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

		columnBeingSorted, sortOrder := parseUrl(r.URL.String())

		deputyClientsDetails, ariaSorting, err := client.GetDeputyClients(ctx, deputyId, columnBeingSorted, sortOrder)
		if err != nil {
			return err
		}

		vars := listClientsVars{
			Path:                 r.URL.Path,
			XSRFToken:            ctx.XSRFToken,
			DeputyClientsDetails: deputyClientsDetails,
			DeputyDetails:        deputyDetails,
			AriaSorting:          ariaSorting,
		}

		if vars.DeputyDetails.ExecutiveCaseManager.EcmId == defaultPATeam {
			vars.ErrorMessage = "An executive case manager has not been assigned. "
		}
		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func parseUrl(url string) (string, string) {
	urlQuery := strings.Split(url, "?")
	if len(urlQuery) >= 2 {
		sortParams := urlQuery[1]
		sortParamsArray := strings.Split(sortParams, ":")
		columnBeingSorted := sortParamsArray[0]
		sortOrder := sortParamsArray[1]
		return columnBeingSorted, sortOrder
	}
	return "", ""
}
