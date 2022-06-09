package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ManageAssuranceVisitsVars struct {
	Path                      string
	XSRFToken                 string
	DeputyDetails             sirius.DeputyDetails
	Error                     string
	Errors                    sirius.ValidationErrors
}

func renderTemplateForAssuranceVisits(tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		vars := ManageAssuranceVisitsVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		}
		fmt.Println(deputyId);
		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

