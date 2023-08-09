package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type FirmInformation interface {
	AddFirmDetails(sirius.Context, sirius.FirmDetails) (int, error)
	AssignDeputyToFirm(sirius.Context, int, int) error
}

type addFirmVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	DeputyId      int
}

func renderTemplateForAddFirm(client FirmInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:
			vars := addFirmVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyId:      deputyId,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			addFirmDetailForm := sirius.FirmDetails{
				FirmName:     r.PostFormValue("name"),
				AddressLine1: r.PostFormValue("address-line-1"),
				AddressLine2: r.PostFormValue("address-line-2"),
				AddressLine3: r.PostFormValue("address-line-3"),
				Town:         r.PostFormValue("town"),
				County:       r.PostFormValue("county"),
				Postcode:     r.PostFormValue("postcode"),
				PhoneNumber:  r.PostFormValue("telephone"),
				Email:        r.PostFormValue("email"),
			}

			firmId, err := client.AddFirmDetails(ctx, addFirmDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := addFirmVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, deputyId, firmId)
			if assignDeputyToFirmErr != nil {
				return assignDeputyToFirmErr
			}

			return Redirect(fmt.Sprintf("/%d?success=newFirm", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
