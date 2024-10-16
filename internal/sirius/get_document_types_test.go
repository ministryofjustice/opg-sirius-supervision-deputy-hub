package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_filterDocTypeByDeputyType(t *testing.T) {
	documentTypes := []model.RefData{{Handle: "ASSURANCE_VISIT", Label: "Assurance visit", Deprecated: false}, {Handle: "CATCH_UP_CALL", Label: "Catch-up call", Deprecated: false}, {Handle: "COMPLAINTS", Label: "Complaints", Deprecated: false}, {Handle: "CORRESPONDENCE", Label: "Correspondence", Deprecated: false}, {Handle: "GENERAL", Label: "General", Deprecated: false}, {Handle: "INDEMNITY_INSURANCE", Label: "Indemnity insurance", Deprecated: false}, {Handle: "NON_COMPLIANCE", Label: "Non-compliance", Deprecated: false}}
	tests := []struct {
		name          string
		DocumentTypes []model.RefData
		DeputyType    string
		want          []model.RefData
	}{
		{
			name:          "Document list types returns without catch-up call for a PRO Deputy",
			DocumentTypes: documentTypes,
			DeputyType:    "PRO",
			want:          []model.RefData{{Handle: "ASSURANCE_VISIT", Label: "Assurance visit", Deprecated: false}, {Handle: "COMPLAINTS", Label: "Complaints", Deprecated: false}, {Handle: "CORRESPONDENCE", Label: "Correspondence", Deprecated: false}, {Handle: "GENERAL", Label: "General", Deprecated: false}, {Handle: "INDEMNITY_INSURANCE", Label: "Indemnity insurance", Deprecated: false}, {Handle: "NON_COMPLIANCE", Label: "Non-compliance", Deprecated: false}},
		},
		{
			name:          "Document list types returns without Indemnity insurance for a PA Deputy",
			DocumentTypes: documentTypes,
			DeputyType:    "PA",
			want:          []model.RefData{{Handle: "ASSURANCE_VISIT", Label: "Assurance visit", Deprecated: false}, {Handle: "CATCH_UP_CALL", Label: "Catch-up call", Deprecated: false}, {Handle: "COMPLAINTS", Label: "Complaints", Deprecated: false}, {Handle: "CORRESPONDENCE", Label: "Correspondence", Deprecated: false}, {Handle: "GENERAL", Label: "General", Deprecated: false}, {Handle: "NON_COMPLIANCE", Label: "Non-compliance", Deprecated: false}},
		},
		{
			name:          "Document list types returns all with a deputy type",
			DocumentTypes: documentTypes,
			DeputyType:    "",
			want:          []model.RefData(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, filterDocTypeByDeputyType(tt.DocumentTypes, tt.DeputyType), "filterDocTypeByDeputyType(%v, %v)", tt.DocumentTypes, tt.DeputyType)
		})
	}
}
