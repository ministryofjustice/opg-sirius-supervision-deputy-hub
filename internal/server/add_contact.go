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
	Path             string
	XSRFToken        string
	DeputyDetails    sirius.DeputyDetails
	Error            string
	Errors           sirius.ValidationErrors
	DeputyId         int
	ContactName      string
	JobTitle         string
	Email            string
	PhoneNumber      string
	OtherPhoneNumber string
	ContactNotes     string
	IsNamedDeputy    string
	IsMainContact    string
}

func convertStringBoolToNullableBoolPointer(stringBool string) *bool {
	if stringBool == "true" {
		return pointerBool(true)
	} else if stringBool == "false" {
		return pointerBool(false)
	}
	return nil
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
				ContactName:      r.PostFormValue("contact-name"),
				JobTitle:         r.PostFormValue("job-title"),
				Email:            r.PostFormValue("email"),
				PhoneNumber:      r.PostFormValue("phone-number"),
				OtherPhoneNumber: r.PostFormValue("other-phone-number"),
				ContactNotes:     r.PostFormValue("contact-notes"),
				IsNamedDeputy:    convertStringBoolToNullableBoolPointer(r.PostFormValue("is-named-deputy")),
				IsMainContact:    convertStringBoolToNullableBoolPointer(r.PostFormValue("is-main-contact")),
			}

			err := client.AddContact(ctx, deputyId, addContactForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := addContactVars{
					Path:             r.URL.Path,
					XSRFToken:        ctx.XSRFToken,
					Errors:           verr.Errors,
					DeputyDetails:    deputyDetails,
					ContactName:      r.PostFormValue("contact-name"),
					JobTitle:         r.PostFormValue("job-title"),
					Email:            r.PostFormValue("email"),
					PhoneNumber:      r.PostFormValue("phone-number"),
					OtherPhoneNumber: r.PostFormValue("other-phone-number"),
					ContactNotes:     r.PostFormValue("contact-notes"),
					IsNamedDeputy:    r.PostFormValue("is-named-deputy"),
					IsMainContact:    r.PostFormValue("is-main-contact"),
				}

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
