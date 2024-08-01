package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strings"
)

type notesVars struct {
	DeputyNotes    sirius.DeputyNoteCollection
	SuccessMessage string
	AppVars
}

type addNoteVars struct {
	Title string
	Note  string
	AppVars
}

type NotesHandler struct {
	router
}

func (h *NotesHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	v.PageName = "Notes"
	if strings.Contains(v.Path, "add-note") {
		v.PageName = "Add a note"
	}

	switch r.Method {
	case http.MethodGet:

		deputyNotes, err := h.Client().GetDeputyNotes(ctx, v.DeputyId())
		if err != nil {
			return err
		}

		successMessage := ""
		if r.URL.Query().Get("success") == "true" {
			successMessage = "Note added"
		}

		vars := notesVars{
			DeputyNotes:    deputyNotes,
			SuccessMessage: successMessage,
			AppVars:        v,
		}

		return h.execute(w, r, vars)

	case http.MethodPost:
		var vars addNoteVars
		var (
			title = r.PostFormValue("title")
			note  = r.PostFormValue("note")
		)

		err := h.Client().AddNote(ctx, title, note, v.DeputyId(), v.UserDetails.ID, v.DeputyType())

		if verr, ok := err.(sirius.ValidationError); ok {
			vars = addNoteVars{
				Title:   title,
				Note:    note,
				AppVars: v,
			}
			vars.Errors = util.RenameErrors(verr.Errors)

			w.WriteHeader(http.StatusBadRequest)
			return h.execute(w, r, vars)
		} else if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d/notes?success=true", v.DeputyId()))
	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
