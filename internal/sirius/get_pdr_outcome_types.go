package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *Client) GetPdrOutcomeTypes(ctx Context) ([]model.RefData, error) {
	pdrOutcomeTypes, err := c.getRefData(ctx, "/pdrOutcome")
	return pdrOutcomeTypes, err
}
