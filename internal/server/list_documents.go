package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"net/http"
)

type DocumentsClient interface {
	GetDeputyDocuments(ctx sirius.Context, deputyId int, sort string) (sirius.DocumentList, error)
}

type DocumentsVars struct {
	DocumentList   sirius.DocumentList
	SuccessMessage string
	AppVars
	Sort string
}

func renderTemplateForDocuments(client DocumentsClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		app.PageName = "Documents"

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addDocument":
			filename := r.URL.Query().Get("filename")
			successMessage = fmt.Sprintf("Document %s added", filename)
		}

		ctx := getContext(r)
		urlParams := r.URL.Query()
		sort := urlbuilder.CreateSortFromURL(urlParams, []string{"receiveddatetime"})

		documentList, err := client.GetDeputyDocuments(ctx, app.DeputyId(), fmt.Sprintf("%s:%s", sort.OrderBy, "desc"))
		if err != nil {
			return err
		}

		vars := DocumentsVars{
			AppVars:        app,
			DocumentList:   documentList,
			SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
