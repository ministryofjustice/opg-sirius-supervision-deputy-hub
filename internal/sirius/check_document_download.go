package sirius

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/telemetry"
)

func (c *Client) CheckDocumentDownload(ctx Context, documentId int) error {
	logger := telemetry.NewLogger("opg-sirius-supervision-deputy-hub ")

	req, err := c.newRequest(ctx, http.MethodHead, fmt.Sprintf(SupervisionAPIPath+"/v1/documents/%d/download", documentId), nil)

	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Error(cerr.Error(), "error", cerr)
		}
	}()

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
