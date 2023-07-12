package server

import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ContactInformation interface {
	AddContact(sirius.Context, int, sirius.Contact) (error)
}

type addContactVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	DeputyId      int
	Form          sirius.Contact
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
			addContactForm := sirius.Contact{
				ContactName:      r.PostFormValue("name"),
				JobTitle:         r.PostFormValue("job-title"),
				Email:            r.PostFormValue("email"),
				PhoneNumber:      r.PostFormValue("phone"),
				OtherPhoneNumber: r.PostFormValue("other-phone"),
				ContactNotes:     r.PostFormValue("notes"),
				IsNamedDeputy:    r.PostFormValue("is-named-deputy"),
				IsMainContact:    r.PostFormValue("is-main-contact"),
			}

			err := client.AddContact(ctx, deputyId, addContactForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := addContactVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
					DeputyDetails: deputyDetails,
					Form: addContactForm,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if(err != nil) {
				return err	
			}

			return Redirect(fmt.Sprintf("/%d/contacts?success=newContact", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}

	}
}