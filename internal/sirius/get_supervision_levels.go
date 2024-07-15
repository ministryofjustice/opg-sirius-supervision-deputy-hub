package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *Client) GetSupervisionLevels(ctx Context) ([]model.RefData, error) {
	supervisionLevels, err := c.getRefData(ctx, "?filter=supervisionLevel")
	return supervisionLevels, err
}
