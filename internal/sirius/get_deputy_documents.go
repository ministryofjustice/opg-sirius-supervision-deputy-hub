package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type DocumentList struct {
	Documents      []model.Document
	Pages          Page
	TotalDocuments int
	Metadata       Metadata
}

func (c *Client) GetDeputyDocuments(ctx Context, deputyId int) (DocumentList, error) {
	var documentList DocumentList
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/persons/%d/documents", deputyId), nil)

	if err != nil {
		return documentList, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return documentList, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return documentList, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return documentList, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&documentList); err != nil {
		return documentList, err
	}

	documentList.formatDocuments()

	fmt.Println("are the document dates updated")
	fmt.Println(documentList.Documents)
	return documentList, err

}

func (d *DocumentList) formatDocuments() {
	var list []model.Document
	for _, existingDocument := range d.Documents {
		existingDocument.CreatedDate = FormatDateTime(SiriusDateTime, existingDocument.CreatedDate, SiriusDate)
		list = append(list, existingDocument)
	}
	d.Documents = list
}
