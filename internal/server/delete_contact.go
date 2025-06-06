package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeleteContact interface {
	GetContactById(ctx sirius.Context, deputyId int, contactId int) (sirius.Contact, error)
	DeleteContact(sirius.Context, int, int) error
}

// Could just use ErrorVars?
type DeleteContactVars struct {
	ContactName string
	AppVars
}

func renderTemplateForDeleteContact(client DeleteContact, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		deputyId, _ := strconv.Atoi(r.PathValue("id"))
		contactId, _ := strconv.Atoi(r.PathValue("contactId"))

		appVars.PageName = "Delete contact"
		vars := DeleteContactVars{
			AppVars: appVars,
		}

		contact, err := client.GetContactById(ctx, deputyId, contactId)
		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			vars.ContactName = contact.ContactName

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			err := client.DeleteContact(ctx, deputyId, contactId)

			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/contacts?success=deletedContact&contactName=%s", deputyId, contact.ContactName))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
