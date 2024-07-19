package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

type AccommodationTypeList struct {
	ClientAccommodation []struct {
		Handle     string `json:"handle"`
		Label      string `json:"label"`
		Deprecated bool   `json:"deprecated"`
	} `json:"clientAccommodation"`
}

func (c *ApiClient) GetAccommodationTypes(ctx Context) ([]model.RefData, error) {
	accommodationTypes, err := c.getRefData(ctx, "?filter=clientAccommodation")
	return accommodationTypes, err
}
