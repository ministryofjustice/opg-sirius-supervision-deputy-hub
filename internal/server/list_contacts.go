package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type DeputyHubContactInformation interface {
	GetDeputyContacts(sirius.Context, int) (sirius.ContactList, error)
}

type ListContactsVars struct {
	SuccessMessage string
	ContactList    sirius.ContactList
	AppVars
}

func renderTemplateForContactTab(client DeputyHubContactInformation, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		contactList, err := client.GetDeputyContacts(ctx, appVars.DeputyId())
		if err != nil {
			return err
		}

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "newContact":
			successMessage = "Contact added"
		case "updatedContact":
			contactName := r.URL.Query().Get("contactName")
			if contactName == "" {
				successMessage = ""
			} else {
				successMessage = contactName + "'s details updated"
			}
		default:
			successMessage = ""
		}

		vars := ListContactsVars{
			AppVars:        appVars,
			ContactList:    contactList,
			SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
