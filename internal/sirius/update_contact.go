package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) UpdateContact(ctx Context, deputyId int, contactId int, manageContactForm ContactForm) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(manageContactForm)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/contacts/%d", deputyId, contactId)

	req, err := c.newRequest(ctx, http.MethodPut, url, &body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return err
}
