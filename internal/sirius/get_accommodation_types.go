package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type AccommodationTypeList struct {
	ClientAccommodation []struct {
		Handle string `json:"handle"`
		Label  string `json:"label"`
	} `json:"clientAccommodation"`
}

func (c *Client) GetAccommodationTypes(ctx Context) ([]model.RefData, error) {
	endpoint := "/api/v1/reference-data?filter=clientAccommodation"

	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v AccommodationTypeList
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	var accommodationTypes []model.RefData
	for _, u := range v.ClientAccommodation {
		accommodationTypes = append(accommodationTypes, u)
	}

	return accommodationTypes, nil
}
