package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ListContactsVars struct {
	Path                 string
	XSRFToken            string
	DeputyDetails        sirius.DeputyDetails
	SuccessMessage       string
	Error                string
}

func renderTemplateForContactTab(client DeputyHubClientInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "newContact":
			successMessage = "Contact added"
		default:
			successMessage = ""
		}

		vars := ListContactsVars{
			Path:                 r.URL.Path,
			XSRFToken:            ctx.XSRFToken,
			DeputyDetails:        deputyDetails,
			SuccessMessage:       successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}