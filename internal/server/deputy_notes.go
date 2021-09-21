package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type DeputyHubNotesInformation interface {
	GetDeputyNotes(sirius.Context, int) (error)
}

type deputyHubNotesVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	DeputyNotes sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
}

func renderTemplateForDeputyHubNotes(client DeputyHubInformation, tmpl Template) Handler {
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
		deputyNotes, err := client.GetDeputyNotes(ctx, deputyId)
		if err != nil {
			return err
		}

		vars := deputyHubVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
			DeputyNotes: deputyNotes,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
