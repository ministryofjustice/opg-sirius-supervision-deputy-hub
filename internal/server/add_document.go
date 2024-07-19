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

type AddDocument interface {
	AddDocument(ctx sirius.Context, file multipart.File, filename string, documentType string, direction string, date string, notes string, deputyId int) error
	GetDocumentDirections(ctx sirius.Context) ([]model.RefData, error)
	GetDocumentTypes(ctx sirius.Context) ([]model.RefData, error)
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

type AddDocumentHandler struct {
	router
}

func (h *AddDocumentHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	v.PageName = "Add a document"

	vars := AddDocumentVars{
		AppVars: v,
		Date:    time.Now().Format("2006-01-02"),
	}

	group, groupCtx := errgroup.WithContext(ctx.Context)

	group.Go(func() error {
		documentDirectionRefData, err := h.Client().GetDocumentDirections(ctx.With(groupCtx))
		if err != nil {
			return err
		}
		vars.DocumentDirectionRefData = documentDirectionRefData
		return nil
	})

	group.Go(func() error {
		documentTypes, err := h.Client().GetDocumentTypes(ctx.With(groupCtx))
		if err != nil {
			return err
		}
		vars.DocumentTypes = documentTypes
		return nil
	})

	if err := group.Wait(); err != nil {
		return err
	}

	switch r.Method {
	case http.MethodGet:
		//implement me!
		return h.execute(w, r, vars, vars.AppVars)

	case http.MethodPost:
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

		//if len(vars.Errors) > 0 {
		//	return tmpl.ExecuteTemplate(w, "page", vars)
		//}

		err = h.Client().AddDocument(ctx, file, handler.Filename, documentType, direction, date, notes, vars.DeputyDetails.ID)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			return h.execute(w, r, vars, v)
		}

		if err != nil {
			return err
		}
		return Redirect(fmt.Sprintf("/%d/documents?success=addDocument&filename=%s", v.DeputyId(), handler.Filename))

	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
