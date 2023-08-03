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
	Firms          []sirius.FirmForList
	Success        bool
	SuccessMessage string
	AppVars
}

func renderTemplateForChangeFirm(client DeputyChangeFirmInformation, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		firms, err := client.GetFirms(ctx)
		if err != nil {
			return err
		}

		vars := changeFirmVars{
			Firms:   firms,
			AppVars: appVars,
		}

		switch r.Method {
		case http.MethodGet:
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
				vars.Errors = verr.Errors
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}
			return Redirect(fmt.Sprintf("/%d?success=firm", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
