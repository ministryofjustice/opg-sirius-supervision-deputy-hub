package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
)

type ManageContact interface {
	GetContactById(ctx sirius.Context, deputyId int, contactId int) (sirius.Contact, error)
	AddContact(sirius.Context, int, sirius.ContactForm) error
	UpdateContact(sirius.Context, int, int, sirius.ContactForm) error
}

type ManageContactVars struct {
	ContactId                     int
	ContactName                   string
	JobTitle                      string
	Email                         string
	PhoneNumber                   string
	OtherPhoneNumber              string
	ContactNotes                  string
	IsNamedDeputy                 string
	IsMainContact                 string
	IsNewContact                  bool
	IsMonthlySpreadsheetRecipient string
	AppVars
}

func renderTemplateForManageContact(client ManageContact, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		deputyId, _ := strconv.Atoi(r.PathValue("id"))
		contactId, _ := strconv.Atoi(r.PathValue("contactId"))

		appVars.PageName = "Add new contact"
		if contactId != 0 {
			appVars.PageName = "Manage contact"
		}

		vars := ManageContactVars{
			AppVars:      appVars,
			IsNewContact: contactId == 0,
		}

		switch r.Method {
		case http.MethodGet:
			if contactId != 0 {
				contact, err := client.GetContactById(ctx, deputyId, contactId)

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
				vars.IsMonthlySpreadsheetRecipient = strconv.FormatBool(contact.IsMonthlySpreadsheetRecipient)
			}
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var successVar string
			var err error

			manageContactForm := sirius.ContactForm{
				ContactName:                   r.PostFormValue("contact-name"),
				JobTitle:                      r.PostFormValue("job-title"),
				Email:                         r.PostFormValue("email"),
				PhoneNumber:                   r.PostFormValue("phone-number"),
				OtherPhoneNumber:              r.PostFormValue("other-phone-number"),
				ContactNotes:                  r.PostFormValue("contact-notes"),
				IsNamedDeputy:                 r.PostFormValue("is-named-deputy"),
				IsMainContact:                 r.PostFormValue("is-main-contact"),
				IsMonthlySpreadsheetRecipient: r.PostFormValue("is-monthly-spreadsheet-recipient"),
			}

			if contactId == 0 {
				err = client.AddContact(ctx, deputyId, manageContactForm)
				successVar = "newContact"
			} else {
				err = client.UpdateContact(ctx, deputyId, contactId, manageContactForm)
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
				vars.IsMonthlySpreadsheetRecipient = manageContactForm.IsMonthlySpreadsheetRecipient

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
