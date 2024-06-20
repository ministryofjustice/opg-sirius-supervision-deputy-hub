package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"io"
	"net/http"
	"strings"
)

func (c *Client) GetRefData(ctx Context, refDataType string) ([]model.RefData, error) {
	var v []model.RefData
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/reference-data%s", refDataType), nil)
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

	if strings.Contains(refDataType, "?filter=") {
		refData, err := unmarshalFilteredRefData(resp.Body, strings.Replace(refDataType, "?filter=", "", -1))

		if err != nil {
			return v, err
		}

		return refData, err

	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	return v, err
}

func unmarshalFilteredRefData(body io.ReadCloser, filter string) ([]model.RefData, error) {
	var refData []model.RefData
	var err error
	var result interface{}

	err = json.NewDecoder(body).Decode(&result)
	if err != nil {
		fmt.Println("Error:", err)
		return refData, err
	}

	dataMap, ok := result.(map[string]interface{})
	if !ok {
		fmt.Println("Invalid JSON structure")
		return refData, err
	}

	data, dataExists := dataMap[filter].([]interface{})

	if !dataExists {
		fmt.Println(fmt.Sprintf("Cannot find key %s in json\n", filter))
		return refData, err
	}

	jsonbody, _ := json.Marshal(data)
	if err != nil {
		return refData, err
	}

	if err := json.Unmarshal(jsonbody, &refData); err != nil {
		return refData, err
	}

	return refData, err
}
