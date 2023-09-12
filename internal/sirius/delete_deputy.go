package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) DeleteDeputy(ctx Context, deputyId int) error {
	var body bytes.Buffer

	url := fmt.Sprintf("/supervision-api/v1/deputies/%d", deputyId)

	req, err := c.newRequest(ctx, http.MethodDelete, url, &body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{Errors: v.ValidationErrors}
		}
		return newStatusError(resp)
	}

	return err
}
