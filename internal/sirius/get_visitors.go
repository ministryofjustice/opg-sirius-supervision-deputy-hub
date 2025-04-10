package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

func (c *Client) GetVisitors(ctx Context) ([]model.Visitor, error) {
	var v []model.Visitor

	req, err := c.newRequest(ctx, http.MethodGet, SupervisionAPIPath + "/v1/visitors", nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&v)

	return v, err
}
