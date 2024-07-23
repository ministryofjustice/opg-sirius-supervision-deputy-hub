package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"io"
	"net/http"
	"os"
)

func (c *Client) GetDocumentById(ctx Context, documentId int) (model.Document, error) {
	var document model.Document

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/documents/%d", documentId), nil)

	if err != nil {
		return document, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return document, err
	}

	io.Copy(os.Stdout, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return document, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return document, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&document)

	return document, err
}
