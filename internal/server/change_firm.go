package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type DeputyChangeFirmInformation interface {
	GetFirms(sirius.Context) ([]sirius.FirmForList, error)
	AssignDeputyToFirm(sirius.Context, int, int) error
}

type changeFirmVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	FirmDetails    []sirius.FirmForList
	Error          string
	Errors         sirius.ValidationErrors
	Success        bool
	SuccessMessage string
}

func renderTemplateForChangeFirm(client DeputyChangeFirmInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		firmDetails, err := client.GetFirms(ctx)
		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:

			if err != nil {
				return err
			}

			vars := changeFirmVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
				FirmDetails:   firmDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var vars changeFirmVars
			newFirm := r.PostFormValue("select-firm")
			AssignToExistingFirmStringIdValue := r.PostFormValue("select-existing-firm")

			if newFirm == "new-firm" {
				return Redirect(fmt.Sprintf("/%d/add-firm", deputyId))
			}

			AssignToFirmId := 0
			if AssignToExistingFirmStringIdValue != "" {
				AssignToFirmId, err = strconv.Atoi(AssignToExistingFirmStringIdValue)
				if err != nil {
					return err
				}
			}

			assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, deputyId, AssignToFirmId)

			if verr, ok := assignDeputyToFirmErr.(sirius.ValidationError); ok {
				vars = changeFirmVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}

				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}
			return Redirect(fmt.Sprintf("/%d?success=firm", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
