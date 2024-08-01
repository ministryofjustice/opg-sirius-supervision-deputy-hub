package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *ApiClient) GetDocumentTypes(ctx Context) ([]model.RefData, error) {
	documentTypes, err := c.getRefData(ctx, "?filter=noteType:deputy")
	return documentTypes, err
}
