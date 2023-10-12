package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type DeleteDeputy interface {
	DeleteDeputy(sirius.Context, int) error
}

type DeleteDeputyVars struct {
	SuccessMessage string
	AppVars
}

func renderTemplateForDeleteDeputy(client DeleteDeputy, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		appVars.PageName = "Delete deputy"

		vars := DeleteDeputyVars{
			AppVars: appVars,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			err := client.DeleteDeputy(ctx, deputyId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			vars.SuccessMessage = fmt.Sprintf("%s %d has been deleted.", vars.DeputyDetails.DisplayName, vars.DeputyDetails.DeputyNumber)

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
