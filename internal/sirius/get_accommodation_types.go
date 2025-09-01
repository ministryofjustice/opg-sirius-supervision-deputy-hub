package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *Client) GetAccommodationTypes(ctx Context) ([]model.RefData, error) {
	accommodationTypes, err := c.getRefData(ctx, "?filter=clientAccommodation")

	if err == nil {
		accommodationTypes = append(
			[]model.RefData{
				{Handle: "HIGH RISK LIVING", Label: "High Risk Living"},
			}, accommodationTypes...)
	}

	return accommodationTypes, err
}
