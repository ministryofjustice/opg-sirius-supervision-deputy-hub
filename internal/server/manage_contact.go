package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strconv"
)

type ManageContactVars struct {
	ContactId        int
	ContactName      string
	JobTitle         string
	Email            string
	PhoneNumber      string
	OtherPhoneNumber string
	ContactNotes     string
	IsNamedDeputy    string
	IsMainContact    string
	IsNewContact     bool
	AppVars
}

type ManageContactsHandler struct {
	router
}

func (h *ManageContactsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	deputyId, _ := strconv.Atoi(r.PathValue("deputyId"))
	contactId, _ := strconv.Atoi(r.PathValue("contactId"))

	v.PageName = "Add new contact"
	if contactId != 0 {
		v.PageName = "Manage contact"
	}

	vars := ManageContactVars{
		AppVars:      v,
		IsNewContact: contactId == 0,
	}

	switch r.Method {
	case http.MethodGet:
		if contactId != 0 {
			contact, err := h.Client().GetContactById(ctx, deputyId, contactId)

			if err != nil {
				return err
			}

			vars.ContactId = contactId
			vars.ContactName = contact.ContactName
			vars.JobTitle = contact.JobTitle
			vars.Email = contact.Email
			vars.PhoneNumber = contact.PhoneNumber
			vars.OtherPhoneNumber = contact.OtherPhoneNumber
			vars.ContactNotes = contact.ContactNotes
			vars.IsNamedDeputy = strconv.FormatBool(contact.IsNamedDeputy)
			vars.IsMainContact = strconv.FormatBool(contact.IsMainContact)
		}
		return h.execute(w, r, vars)

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

		if contactId == 0 {
			err = h.Client().AddContact(ctx, deputyId, manageContactForm)
			successVar = "newContact"
		} else {
			err = h.Client().UpdateContact(ctx, deputyId, contactId, manageContactForm)
			successVar = "updatedContact&contactName=" + r.PostFormValue("contact-name")
		}

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			vars.ContactName = manageContactForm.ContactName
			vars.JobTitle = manageContactForm.JobTitle
			vars.Email = manageContactForm.Email
			vars.PhoneNumber = manageContactForm.PhoneNumber
			vars.OtherPhoneNumber = manageContactForm.OtherPhoneNumber
			vars.ContactNotes = manageContactForm.ContactNotes
			vars.IsNamedDeputy = manageContactForm.IsNamedDeputy
			vars.IsMainContact = manageContactForm.IsMainContact
			vars.IsNewContact = contactId == 0

			return h.execute(w, r, vars)
		}

		if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d/contacts?success=%s", deputyId, successVar))
	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
