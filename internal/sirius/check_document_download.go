package sirius

import (
	"fmt"
	"net/http"
)

func (c *Client) CheckDocumentDownload(ctx Context, documentId int) error {
	req, err := c.newRequest(ctx, http.MethodHead, fmt.Sprintf(SupervisionAPIPath+"/v1/documents/%d/download", documentId), nil)

	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode == http.StatusBadRequest {
		return newStatusError(resp) // or create a specific error for this case
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil // Success case
}
