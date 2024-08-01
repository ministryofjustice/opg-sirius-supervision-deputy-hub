package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"net/http"
)

type DocumentsVars struct {
	DocumentList   sirius.DocumentList
	SuccessMessage string
	AppVars
	Sort string
}

type ListDocumentsHandler struct {
	router
}

func (h *ListDocumentsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	v.PageName = "Documents"

	var successMessage string
	switch r.URL.Query().Get("success") {
	case "addDocument":
		filename := r.URL.Query().Get("filename")
		successMessage = fmt.Sprintf("Document %s added", filename)
	case "replaceDocument":
		previousFilename := r.URL.Query().Get("previousFilename")
		filename := r.URL.Query().Get("filename")
		successMessage = fmt.Sprintf("Document %s has been replaced by %s", previousFilename, filename)
	}

	urlParams := r.URL.Query()
	sort := urlbuilder.CreateSortFromURL(urlParams, []string{"receiveddatetime"})

	documentList, err := h.Client().GetDeputyDocuments(ctx, v.DeputyId(), fmt.Sprintf("%s:%s", sort.OrderBy, "desc"))
	if err != nil {
		return err
	}

	vars := DocumentsVars{
		AppVars:        v,
		DocumentList:   documentList,
		SuccessMessage: successMessage,
	}

	return h.execute(w, r, vars)
}
