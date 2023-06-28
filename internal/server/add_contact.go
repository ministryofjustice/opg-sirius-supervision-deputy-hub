package server

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ContactInformation interface {
	AddContactDetails(sirius.Context, int, sirius.ContactDetails) (error)
}

type addContactVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	DeputyId      int
}

func renderTemplateForAddContact(client ContactInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:
			vars := addContactVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyId:      deputyId,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:

			addContactDetailForm := sirius.ContactDetails{
				Name:           r.PostFormValue("name"),
				JobTitle:       r.PostFormValue("job-title"),
				Email:          r.PostFormValue("email"),
				Phone:          r.PostFormValue("phone"),
				SecondaryPhone: r.PostFormValue("phone-secondary"),
				Notes:          r.PostFormValue("notes"),
				NamedDeputy:    r.PostFormValue("named-deputy") == "yes",
				MainContact:    r.PostFormValue("main-contact") == "yes",
			}

			err := client.AddContactDetails(ctx, deputyId, addContactDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := addFirmVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/%d/contacts?success=newContact", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}

	}
}