package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type DocumentList struct {
	Documents      []model.Document `json:"documents"`
	Pages          Page             `json:"pages"`
	TotalDocuments int              `json:"total"`
	Metadata       Metadata         `json:"metadata"`
}

func (c *Client) GetDeputyDocuments(ctx Context, deputyId int, sort string) (DocumentList, error) {
	var documentList DocumentList
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(SupervisionAPIPath+"/v1/persons/%d/documents?&sort=%s", deputyId, sort), nil)

	if err != nil {
		return documentList, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return documentList, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return documentList, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return documentList, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&documentList); err != nil {
		return documentList, err
	}

	documentList.Documents = formatDocuments(documentList.Documents)

	return documentList, err
}

func formatDocuments(documents []model.Document) []model.Document {
	for key, document := range documents {
		documents[key].CreatedDate = FormatDateTime(SiriusDateTime, document.CreatedDate, SiriusDate)
		documents[key].ReceivedDateTime = FormatDateTime(SiriusDateTime, document.ReceivedDateTime, SiriusDate)
	}
	return documents
}
