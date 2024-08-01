package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *ApiClient) GetDocumentDirections(ctx Context) ([]model.RefData, error) {
	documentDirections, err := c.getRefData(ctx, "/documentDirection")
	return documentDirections, err
}
