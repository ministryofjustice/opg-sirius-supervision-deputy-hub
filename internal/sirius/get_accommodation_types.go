package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *Client) GetAccommodationTypes(ctx Context) ([]model.RefData, error) {
	accommodationTypes, err := c.getRefData(ctx, "?filter=clientAccommodation")
	return accommodationTypes, err
}
