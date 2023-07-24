package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type ManageContact interface {
	GetContactById(ctx sirius.Context, deputyId int, contactId int) (sirius.Contact, error)
	AddContact(sirius.Context, int, sirius.ContactForm) error
	UpdateContact(sirius.Context, int, int, sirius.ContactForm) error
}

type ManageContactVars struct {
	Path             string
	XSRFToken        string
	DeputyDetails    sirius.DeputyDetails
	Error            string
	Errors           sirius.ValidationErrors
	Contact          sirius.Contact
	ErrorNote        string
	ContactName      string
	JobTitle         string
	Email            string
	PhoneNumber      string
	OtherPhoneNumber string
	ContactNotes     string
	IsNamedDeputy    string
	IsMainContact    string
	IsNewContact     bool
}

func renderTemplateForManageContact(client ManageContact, tmpl Template, isNewContact bool) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		var contactId int
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		if !isNewContact {
			contactId, _ = strconv.Atoi(routeVars["contactId"])
		}

		vars := ManageContactVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
			IsNewContact:  isNewContact,
		}

		switch r.Method {
		case http.MethodGet:
			if !isNewContact {
				contact, err := client.GetContactById(ctx, deputyId, contactId)
				if err != nil {
					return err
				}

				vars.ContactName = contact.ContactName
				vars.JobTitle = contact.JobTitle
				vars.Email = contact.Email
				vars.PhoneNumber = contact.PhoneNumber
				vars.OtherPhoneNumber = contact.OtherPhoneNumber
				vars.ContactNotes = contact.ContactNotes
				vars.IsNamedDeputy = strconv.FormatBool(contact.IsNamedDeputy)
				vars.IsMainContact = strconv.FormatBool(contact.IsMainContact)
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var successVar string
			var err error

			manageContactForm := sirius.ContactForm{
				ContactName:      r.PostFormValue("contact-name"),
				JobTitle:         r.PostFormValue("job-title"),
				Email:            r.PostFormValue("email"),
				PhoneNumber:      r.PostFormValue("phone-number"),
				OtherPhoneNumber: r.PostFormValue("other-phone-number"),
				ContactNotes:     r.PostFormValue("contact-notes"),
				IsNamedDeputy:    r.PostFormValue("is-named-deputy"),
				IsMainContact:    r.PostFormValue("is-main-contact"),
			}

			if isNewContact {
				err = client.AddContact(ctx, deputyId, manageContactForm)
				successVar = "newContact"
			} else {
				err = client.UpdateContact(ctx, deputyId, contactId, manageContactForm)
				successVar = "updatedContact"
			}

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := ManageContactVars{
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

			return Redirect(fmt.Sprintf("/%d/contacts?success=%s", deputyId, successVar))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
