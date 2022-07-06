package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type AssuranceVisit interface {
	GetAssuranceVisits(ctx sirius.Context, deputyId int) ([]sirius.AssuranceVisits, error)
}

type AssuranceVisitsVars struct {
	Path            string
	XSRFToken       string
	DeputyDetails   sirius.DeputyDetails
	Error           string
	Errors          sirius.ValidationErrors
	SuccessMessage  string
	AssuranceVisits []sirius.AssuranceVisits
}

func renderTemplateForAssuranceVisits(client AssuranceVisit, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addAssuranceVisit":
			successMessage = "Assurance process updated"
		case "manageAssuranceVisit":
			successMessage = "Assurance visit updated"
		default:
			successMessage = ""
		}

		visits, err := client.GetAssuranceVisits(ctx, deputyId)
		if err != nil {
			return err
		}

		vars := AssuranceVisitsVars{
			Path:            r.URL.Path,
			XSRFToken:       ctx.XSRFToken,
			DeputyDetails:   deputyDetails,
			SuccessMessage:  successMessage,
			AssuranceVisits: visits,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
