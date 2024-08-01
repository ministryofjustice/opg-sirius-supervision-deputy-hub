package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *ApiClient) GetVisitOutcomeTypes(ctx Context) ([]model.RefData, error) {
	visitOutcomeTypes, err := c.getRefData(ctx, "/visitOutcome")
	return visitOutcomeTypes, err
}
