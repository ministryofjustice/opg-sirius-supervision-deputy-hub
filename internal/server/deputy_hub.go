package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubInformation interface {
	GetDeputyDetails(sirius.Context) (sirius.DeputyDetails, error)
}

type deputyHubVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
}

func renderTemplateForDeputyHub(client DeputyHubInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		//deputyDetails, err := client.GetDeputyDetails(ctx)
		//if err != nil {
		//	return err
		//}

		vars := deputyHubVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			//DeputyDetails: deputyDetails,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
