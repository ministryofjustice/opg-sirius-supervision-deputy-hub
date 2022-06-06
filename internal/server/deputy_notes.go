package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
	"strings"
)

type DeputyHubNotesInformation interface {
	GetDeputyNotes(sirius.Context, int) (sirius.DeputyNoteCollection, error)
	AddNote(ctx sirius.Context, title, note string, deputyId, userId int, deputyType string) error
	GetUserDetails(sirius.Context) (sirius.UserDetails, error)
}

type deputyHubNotesVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	DeputyNotes    sirius.DeputyNoteCollection
	Error          string
	Errors         sirius.ValidationErrors
	ErrorMessage   string
	Success        bool
	SuccessMessage string
}

type addNoteVars struct {
	Path          string
	XSRFToken     string
	Title         string
	Note          string
	Success       bool
	Error         sirius.ValidationError
	Errors        sirius.ValidationErrors
	ErrorMessage  string
	DeputyDetails sirius.DeputyDetails
}

func hasSuccessInUrl(url string, prefix string) bool {
	urlTrim := strings.TrimPrefix(url, prefix)
	return urlTrim == "?success=true"
}

func renderTemplateForDeputyHubNotes(client DeputyHubNotesInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:

			deputyNotes, err := client.GetDeputyNotes(ctx, deputyId)
			if err != nil {
				return err
			}

			hasSuccess := hasSuccessInUrl(r.URL.String(), "/"+strconv.Itoa(deputyId)+"/notes")

			vars := deputyHubNotesVars{
				Path:           r.URL.Path,
				XSRFToken:      ctx.XSRFToken,
				DeputyDetails:  deputyDetails,
				DeputyNotes:    deputyNotes,
				Success:        hasSuccess,
				SuccessMessage: "Note added",
			}

			vars.ErrorMessage = checkForDefaultEcmId(deputyDetails.ExecutiveCaseManager.EcmId, defaultPATeam)
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var vars addNoteVars
			var (
				title = r.PostFormValue("title")
				note  = r.PostFormValue("note")
			)

			userId, err := client.GetUserDetails(ctx)
			if err != nil {
				return err
			}

			err = client.AddNote(ctx, title, note, deputyId, userId.ID, deputyDetails.DeputyType.Handle)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars = addNoteVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					Title:         title,
					Note:          note,
					Errors:        verr.Errors,
					DeputyDetails: deputyDetails,
				}

				vars.ErrorMessage = checkForDefaultEcmId(deputyDetails.ExecutiveCaseManager.EcmId, defaultPATeam)

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
