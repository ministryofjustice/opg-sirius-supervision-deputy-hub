package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type AccommodationTypeList struct {
	ClientAccommodation []struct {
		Handle     string `json:"handle"`
		Label      string `json:"label"`
		Deprecated bool   `json:"deprecated"`
	} `json:"clientAccommodation"`
}

func (c *ApiClient) GetAccommodationTypes(ctx Context) ([]model.RefData, error) {
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
		if !u.Deprecated {
			accommodationTypes = append(accommodationTypes, u)
		}
	}

	accommodationTypes = append(
		[]model.RefData{
			{Handle: "HIGH RISK LIVING", Label: "High Risk Living"},
		}, accommodationTypes...)

	return accommodationTypes, nil
}
