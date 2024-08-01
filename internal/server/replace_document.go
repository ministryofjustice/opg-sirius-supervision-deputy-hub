package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"golang.org/x/sync/errgroup"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type ReplaceDocument interface {
	ReplaceDocument(ctx sirius.Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId, documentId int) error
	GetDocumentDirections(ctx sirius.Context) ([]model.RefData, error)
	GetDocumentTypes(ctx sirius.Context) ([]model.RefData, error)
	GetDocumentById(ctx sirius.Context, deputyId, documentId int) (model.Document, error)
}

type ReplaceDocumentVars struct {
	SuccessMessage   string
	OriginalDocument model.Document
	AppVars
	DocumentDirectionRefData []model.RefData
	DocumentTypes            []model.RefData
	DocumentType             string
	Direction                string
	Date                     string
	Notes                    string
}

type ReplaceDocumentHandler struct {
	router
}

func (h *ReplaceDocumentHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	v.PageName = "Replace a document"

	routeVars := mux.Vars(r)
	documentId, _ := strconv.Atoi(routeVars["documentId"])

	vars := ReplaceDocumentVars{
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

	group.Go(func() error {
		originalDocument, err := h.Client().GetDocumentById(ctx.With(groupCtx), vars.DeputyDetails.ID, documentId)
		if err != nil {
			return err
		}
		if (originalDocument != model.Document{}) {
			vars.OriginalDocument = originalDocument
			newTime, err := time.Parse("02/01/2006 15:04:05", originalDocument.ReceivedDateTime)
			if err != nil {
				return err
			}
			vars.OriginalDocument.ReformattedTime = newTime.Format("02/01/2006")
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return err
	}

	switch r.Method {
	case http.MethodGet:
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

		if len(vars.Errors) > 0 {
			return h.execute(w, r, vars, v)
		}

		ctx := getContext(r)
		err = h.Client().ReplaceDocument(ctx, file, handler.Filename, documentType, direction, date, notes, vars.DeputyDetails.ID, vars.OriginalDocument.Id)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			return h.execute(w, r, vars, v)
		}

		if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d/documents?success=replaceDocument&previousFilename=%s&filename=%s", v.DeputyId(), vars.OriginalDocument.FriendlyDescription, handler.Filename))
	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
