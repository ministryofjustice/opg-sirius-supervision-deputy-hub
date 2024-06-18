package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type DocumentsClient interface {
	GetDeputyDocuments(ctx sirius.Context, deputyId int) (sirius.DocumentList, error)
}

type DocumentsVars struct {
	DocumentList   sirius.DocumentList
	SuccessMessage string
	AppVars
}

func renderTemplateForDocuments(client DocumentsClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)

		documentList, err := client.GetDeputyDocuments(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		vars := DocumentsVars{
			AppVars:      app,
			DocumentList: documentList,
			//SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
