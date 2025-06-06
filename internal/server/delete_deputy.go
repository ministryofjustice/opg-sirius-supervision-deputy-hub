package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
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
		deputyId, _ := strconv.Atoi(r.PathValue("id"))

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
				vars.Errors = util.RenameErrors(verr.Errors)

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
