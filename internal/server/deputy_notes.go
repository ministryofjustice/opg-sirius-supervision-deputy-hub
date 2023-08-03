package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
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
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:

			deputyNotes, err := client.GetDeputyNotes(ctx, deputyId)
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
				AppVars:        appVars,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var vars addNoteVars
			var (
				title = r.PostFormValue("title")
				note  = r.PostFormValue("note")
			)

			err := client.AddNote(ctx, title, note, deputyId, appVars.UserDetails.ID, appVars.DeputyDetails.DeputyType.Handle)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars = addNoteVars{
					Title:   title,
					Note:    note,
					AppVars: appVars,
				}
				vars.Errors = verr.Errors

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/notes?success=true", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
