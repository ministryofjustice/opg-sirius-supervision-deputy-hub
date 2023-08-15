package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type DeputyHubNotesInformation interface {
	GetDeputyNotes(sirius.Context, int) (sirius.DeputyNoteCollection, error)
	AddNote(ctx sirius.Context, title, note string, deputyId, userId int, deputyType string) error
}

type deputyHubNotesVars struct {
	DeputyNotes    sirius.DeputyNoteCollection
	SuccessMessage string
	AppVars
}

type addNoteVars struct {
	Title string
	Note  string
	AppVars
}

func renderTemplateForDeputyHubNotes(client DeputyHubNotesInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		switch r.Method {
		case http.MethodGet:

			deputyNotes, err := client.GetDeputyNotes(ctx, app.DeputyId())
			if err != nil {
				return err
			}

			successMessage := ""
			if r.URL.Query().Get("success") == "true" {
				successMessage = "Note added"
			}

			vars := deputyHubNotesVars{
				DeputyNotes:    deputyNotes,
				SuccessMessage: successMessage,
				AppVars:        app,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var vars addNoteVars
			var (
				title = r.PostFormValue("title")
				note  = r.PostFormValue("note")
			)

			err := client.AddNote(ctx, title, note, app.DeputyId(), app.UserDetails.ID, app.DeputyDetails.DeputyType.Handle)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars = addNoteVars{
					Title:   title,
					Note:    note,
					AppVars: app,
				}
				vars.Errors = verr.Errors

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/notes?success=true", app.DeputyId()))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
