package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
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

type DeleteContactHandler struct {
	router
}

func (h *DeleteContactHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	deputyId, _ := strconv.Atoi(r.PathValue("deputyId"))
	contactId, _ := strconv.Atoi(r.PathValue("contactId"))

	v.PageName = "Delete contact"
	vars := DeleteContactVars{
		AppVars: v,
	}

	contact, err := h.Client().GetContactById(ctx, deputyId, contactId)
	if err != nil {
		return err
	}

	switch r.Method {
	case http.MethodGet:
		vars.ContactName = contact.ContactName
		return h.execute(w, r, vars, v)

	case http.MethodPost:
		err := h.Client().DeleteContact(ctx, deputyId, contactId)
		if err != nil {
			return err
		}
		return Redirect(fmt.Sprintf("/%d/contacts?success=deletedContact&contactName=%s", deputyId, contact.ContactName))

	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
