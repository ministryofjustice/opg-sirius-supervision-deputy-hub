package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *Client) GetGCMIssueTypes(ctx Context) ([]model.RefData, error) {
	gcmIssueTypes, err := c.getRefData(ctx, "/gcmIssueType")
	return gcmIssueTypes, err
}
