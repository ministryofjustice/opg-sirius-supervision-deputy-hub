package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"golang.org/x/sync/errgroup"
	"mime/multipart"
	"net/http"
	"time"
)

type AddDocumentClient interface {
	AddDocument(ctx sirius.Context, file multipart.File, filename string, documentType string, direction string, date string, notes string, deputyId int) error
	GetDocumentDirections(ctx sirius.Context) ([]model.RefData, error)
	GetDocumentTypes(ctx sirius.Context, deputyType string) ([]model.RefData, error)
}

type AddDocumentVars struct {
	SuccessMessage string
	AppVars
	DocumentDirectionRefData []model.RefData
	DocumentTypes            []model.RefData
	DocumentType             string
	Direction                string
	Date                     string
	Notes                    string
}

func renderTemplateForAddDocument(client AddDocumentClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		app.PageName = "Add a document"

		vars := AddDocumentVars{
			AppVars: app,
			Date:    time.Now().Format("2006-01-02"),
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			documentDirectionRefData, err := client.GetDocumentDirections(ctx.With(groupCtx))
			if err != nil {
				return err
			}
			vars.DocumentDirectionRefData = documentDirectionRefData
			return nil
		})

		group.Go(func() error {
			documentTypes, err := client.GetDocumentTypes(ctx.With(groupCtx), vars.DeputyType())
			if err != nil {
				return err
			}
			vars.DocumentTypes = documentTypes
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			vars.Errors = sirius.ValidationErrors{}

			// Specify max file size to 100mb
			err := r.ParseMultipartForm(100 << 20)
			if err != nil {
				return err
			}

			file, handler, err := r.FormFile("document-upload")
			if err != nil {
				vars.Errors["document-upload"] = map[string]string{"": "Select a file to attach"}
			}

			documentType := r.PostFormValue("documentType")
			direction := r.PostFormValue("documentDirection")
			date := r.PostFormValue("documentDate")
			notes := r.PostFormValue("notes")

			vars.DocumentType = documentType
			vars.Direction = direction
			vars.Date = date
			vars.Notes = notes

			if len(vars.Errors) > 0 {
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			ctx := getContext(r)
			err = client.AddDocument(ctx, file, handler.Filename, documentType, direction, date, notes, vars.DeputyDetails.ID)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/documents?success=addDocument&filename=%s", app.DeputyId(), handler.Filename))
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}

}
