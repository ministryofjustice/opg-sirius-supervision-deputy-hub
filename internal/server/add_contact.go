package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type ContactInformation interface {
	AddContact(sirius.Context, int, sirius.Contact) error
}

type addContactVars struct {
	ContactName      string
	JobTitle         string
	Email            string
	PhoneNumber      string
	OtherPhoneNumber string
	ContactNotes     string
	IsNamedDeputy    string
	IsMainContact    string
	AppVars
}

func renderTemplateForAddContact(client ContactInformation, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		vars := addContactVars{
			AppVars: appVars,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:
			addContactForm := sirius.Contact{
				ContactName:      r.PostFormValue("contact-name"),
				JobTitle:         r.PostFormValue("job-title"),
				Email:            r.PostFormValue("email"),
				PhoneNumber:      r.PostFormValue("phone-number"),
				OtherPhoneNumber: r.PostFormValue("other-phone-number"),
				ContactNotes:     r.PostFormValue("contact-notes"),
				IsNamedDeputy:    r.PostFormValue("is-named-deputy"),
				IsMainContact:    r.PostFormValue("is-main-contact"),
			}

			err := client.AddContact(ctx, deputyId, addContactForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				vars.ContactName = addContactForm.ContactName
				vars.JobTitle = addContactForm.JobTitle
				vars.Email = addContactForm.Email
				vars.PhoneNumber = addContactForm.PhoneNumber
				vars.OtherPhoneNumber = addContactForm.OtherPhoneNumber
				vars.ContactNotes = addContactForm.ContactNotes
				vars.IsNamedDeputy = addContactForm.IsNamedDeputy
				vars.IsMainContact = addContactForm.IsMainContact

				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/contacts?success=newContact", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}

	}
}
