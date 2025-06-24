package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type CheckDocumentDownload interface {
	CheckDocumentDownload(ctx sirius.Context, documentId int) error
}

func renderTemplateForCheckDocument(client CheckDocumentDownload) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodHead && r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		documentId, _ := strconv.Atoi(r.PathValue("documentId"))

		err := client.CheckDocumentDownload(ctx, documentId)
		if err != nil {
			// Handle the error - could return 400 status or whatever you need
			return err
		}

		w.WriteHeader(http.StatusOK)
		return nil
	}
}
