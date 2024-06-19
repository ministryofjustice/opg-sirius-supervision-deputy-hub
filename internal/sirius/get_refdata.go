package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

func (c *Client) GetRefData(ctx Context, refDataType string) ([]model.RefData, error) {
	var v []model.RefData

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/reference-data/%s", refDataType), nil)
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

	fmt.Print("ref data")
	fmt.Println(v)

	return v, err
}
