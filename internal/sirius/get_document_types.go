package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *Client) GetDocumentTypes(ctx Context, deputyType string) ([]model.RefData, error) {
	documentTypes, err := c.getRefData(ctx, "?filter=noteType:deputy")
	documentTypesFilteredByDeputyType := filterDocTypeByDeputyType(documentTypes, deputyType)
	return documentTypesFilteredByDeputyType, err
}

func filterDocTypeByDeputyType(DocumentTypes []model.RefData, DeputyType string) []model.RefData {
	var deputySpecificDocTypes []model.RefData
	for _, v := range DocumentTypes {
		if v.Handle != "INDEMNITY_INSURANCE" && DeputyType == "PA" {
			deputySpecificDocTypes = append(deputySpecificDocTypes, v)
		}
		if v.Handle != "CATCH_UP_CALL" && DeputyType == "PRO" {
			deputySpecificDocTypes = append(deputySpecificDocTypes, v)
		}
	}
	return deputySpecificDocTypes
}
