package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *Client) GetDocumentTypes(ctx Context, deputyType string) ([]model.RefData, error) {
	documentTypes, err := c.getRefData(ctx, "?filter=noteType:deputy")
	var documentTypesFilteredByDeputyType []model.RefData
	documentTypesFilteredByDeputyType = filterDocTypeByDeputyType(documentTypes, deputyType)
	return documentTypesFilteredByDeputyType, err
}

func filterDocTypeByDeputyType(DocumentTypes []model.RefData, DeputyType string) []model.RefData {
	var deputySpecificDocTypes []model.RefData
	for _, v := range DocumentTypes {
		if v.Label != "Indemnity insurance" && DeputyType == "PA" {
			deputySpecificDocTypes = append(deputySpecificDocTypes, v)
		}
		if v.Label != "Catch-up call" && DeputyType == "PRO" {
			deputySpecificDocTypes = append(deputySpecificDocTypes, v)
		}
	}
	return deputySpecificDocTypes
}
