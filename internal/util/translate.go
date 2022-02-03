package util

var mappings = map[string]string{
	"FIELD.firstname":    "First name",
	"FIELD.surname":      "Surname",
	"FIELD.deputyName":   "Deputy name",
	"FIELD.telephone":    "Telephone",
	"FIELD.email":        "Email address",
	"FIELD.addressLine1": "Address line 1",
	"FIELD.addressLine2": "Address line 2",
	"FIELD.addressLine3": "Address line 3",
	"FIELD.town":         "Town or City",
	"FIELD.postcode":     "Postcode",
	"FIELD.county":       "County",
	"FIELD.country":      "Country",
	"FIELD.annualBillingInvoice":      "Annual billing invoice",
	"FIELD.panelDeputy":               "Panel deputy",
	"FIELD.complaints":                "Complaints",
	"FIELD.otherImportantInformation": "Other important information",
	"VALUE.true":                      "Yes",
	"VALUE.false":                     "No",
}

func Translate(prefix string, s string) string {
	val, ok := mappings[prefix+"."+s]
	if !ok {
		return s
	}
	return val
}
