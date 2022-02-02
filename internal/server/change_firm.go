package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

func renderTemplateForChangeFirm(client DeputyHubInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:
			deputyDetails, err := client.GetDeputyDetails(ctx, defaultPATeam, deputyId)

			if err != nil {
				return err
			}

			vars := deputyHubVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			newFirm := r.PostFormValue("select-firm")

			if newFirm == "new-firm" {
				return Redirect(fmt.Sprintf("/%d/add-firm", deputyId))
			}

			vars := deputyHubVars{
				Path:      r.URL.Path,
				XSRFToken: ctx.XSRFToken,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
