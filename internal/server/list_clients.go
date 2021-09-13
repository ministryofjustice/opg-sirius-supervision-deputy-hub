package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type DeputyHubClientInformation interface {
	GetDeputyClients(sirius.Context, int) (sirius.DeputyClientDetails, error)
	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
}

type listClientsVars struct {
	Path          string
	XSRFToken     string
	DeputyClientsDetails sirius.DeputyClientDetails
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
}

func listClients(client DeputyHubClientInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails, err := client.GetDeputyDetails(ctx, deputyId)
		deputyClientsDetails, err := client.GetDeputyClients(ctx, deputyId)
		if err != nil {
			return err
		}

		vars := listClientsVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyClientsDetails: deputyClientsDetails,
			DeputyDetails: deputyDetails,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}