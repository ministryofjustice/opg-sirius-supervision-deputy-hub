package sirius

import (
	"bytes"
	"fmt"
	"net/http"
)

func (c *Client) DeleteContact(ctx Context, deputyId int, contactId int) error {
	var body bytes.Buffer

	url := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/contacts/%d", deputyId, contactId)

	req, err := c.newRequest(ctx, http.MethodDelete, url, &body)

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
		return newStatusError(resp)
	}

	return err
}
