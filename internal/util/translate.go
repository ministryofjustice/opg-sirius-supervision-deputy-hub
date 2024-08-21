package util

import (
	"time"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

var translationMappings = map[string]string{
	"FIELD.firstname":                 "First name",
	"FIELD.surname":                   "Surname",
	"FIELD.deputyName":                "Deputy name",
	"FIELD.telephone":                 "Telephone",
	"FIELD.email":                     "Email address",
	"FIELD.addressLine1":              "Address line 1",
	"FIELD.addressLine2":              "Address line 2",
	"FIELD.addressLine3":              "Address line 3",
	"FIELD.town":                      "Town or City",
	"FIELD.postcode":                  "Postcode",
	"FIELD.county":                    "County",
	"FIELD.country":                   "Country",
	"FIELD.annualBillingInvoice":      "Annual billing invoice",
	"FIELD.panelDeputy":               "Panel deputy",
	"FIELD.independentVisitorCharges": "Independent visitor charges",
	"FIELD.bankCharges":               "Bank charges",
	"FIELD.apad":                      "APAD",
	"FIELD.reportSystemType":          "Report system",
	"FIELD.reportSystem":              "Report system",
	"FIELD.monthlySpreadsheet":        "Monthly spreadsheet",
	"FIELD.complaints":                "Complaints",
	"FIELD.otherImportantInformation": "Other important information",
	"VALUE.YES":                       "Yes",
	"VALUE.NO":                        "No",
	"VALUE.true":                      "Yes",
	"VALUE.false":                     "No",
	"VALUE.UNKNOWN":                   "Unknown",
}

func Translate(prefix string, s string) string {
	val, ok := translationMappings[prefix+"."+s]
	if !ok {
		return s
	}
	return val
}

type pair struct {
	k string
	v string
}

var currentDate = time.Now().Format("02/01/2006")

var validationMappings = map[string]map[string]pair{
	// firm
	"firmName": {
		"stringLengthTooLong": pair{"firmName", "The firm name must be 255 characters or fewer"},
		"isEmpty":             pair{"firmName", "The firm name is required and can't be empty"},
	},
	"firmId": {
		"notGreaterThanInclusive": pair{"existing-firm", "Enter a firm name or number"},
	},
	// address
	"addressLine1": {
		"stringLengthTooLong": pair{"addressLine1", "The building or street must be 255 characters or fewer"},
	},
	"addressLine2": {
		"stringLengthTooLong": pair{"addressLine2", "Address line 2 must be 255 characters or fewer"},
	},
	"addressLine3": {
		"stringLengthTooLong": pair{"addressLine3", "Address line 3 must be 255 characters or fewer"},
	},
	"town": {
		"stringLengthTooLong": pair{"town", "The town or city must be 255 characters or fewer"},
	},
	"county": {
		"stringLengthTooLong": pair{"county", "The county must be 255 characters or fewer"},
	},
	"postcode": {
		"stringLengthTooLong": pair{"postcode", "The postcode must be 255 characters or fewer"},
	},
	"phoneNumber": {
		"stringLengthTooLong": pair{"phoneNumber", "The telephone number must be 255 characters or fewer"},
		"isEmpty":             pair{"phoneNumber", "Enter a telephone number"},
	},
	"workPhoneNumber": {
		"stringLengthTooLong": pair{"workPhoneNumber", "The telephone number must be 255 characters or fewer"},
	},
	"email": {
		"stringLengthTooLong":         pair{"email", "The email must be 255 characters or fewer"},
		"isEmpty":                     pair{"email", "Enter an email address"},
		"emailAddressInvalidFormat":   pair{"email", "Enter an email address in the correct format, like name@example.com"},
		"emailAddressInvalidHostname": pair{"email", "The email address is invalid"},
		"hostnameInvalidHostname":     pair{"email", "The email address is invalid"},
		"hostnameLocalNameNotAllowed": pair{"email", "The email address is invalid"},
	},
	// note
	"name": {
		"stringLengthTooLong": pair{"1-title", "The title must be 255 characters or fewer"},
		"isEmpty":             pair{"1-title", "Enter a title for the note"},
	},
	"description": {
		"stringLengthTooLong": pair{"2-note", "The note must be 1000 characters or fewer"},
		"isEmpty":             pair{"2-note", "Enter a note"},
	},
	// deputy
	"organisationName": {
		"stringLengthTooLong": pair{"organisationName", "The deputy name must be 255 characters or fewer"},
		"isEmpty":             pair{"organisationName", "Enter a deputy name"},
	},
	"organisationTeamOrDepartmentName": {
		"stringLengthTooLong": pair{"organisationTeamOrDepartmentName", "The team or department must be 255 characters or fewer"},
	},
	"firstname": {
		"stringLengthTooLong": pair{"firstname", "The deputy first name must be 255 characters or fewer"},
		"isEmpty":             pair{"firstname", "The deputy first name is required and can't be empty"},
	},
	"surname": {
		"stringLengthTooLong": pair{"surname", "The deputy surname must be 255 characters or fewer"},
		"isEmpty":             pair{"surname", "The deputy surname is required and can't be empty"},
	},

	// deputy contact
	"contactName": {
		"stringLengthTooLong": pair{"contactName", "The name must be 255 characters or fewer"},
		"isEmpty":             pair{"contactName", "Enter a name"},
	},
	"jobTitle": {
		"stringLengthTooLong": pair{"jobTitle", "The job title must be 255 characters or fewer"},
	},
	"otherPhoneNumber": {
		"stringLengthTooLong": pair{"otherPhoneNumber", "The other telephone number must be 255 characters or fewer"},
	},
	"contactNotes": {
		"stringLengthTooLong": pair{"contactNotes", "The note must be 255 characters or fewer"},
	},
	"isMainContact": {
		"isEmpty": pair{"isMainContact", "Select whether this contact is a main contact"},
	},
	"isNamedDeputy": {
		"isEmpty": pair{"isNamedDeputy", "Select whether this contact is the named deputy"},
	},
	// task
	"taskType": {
		"isEmpty": pair{"taskType", "Select the task type"},
	},
	"dueDate": {
		"isEmpty": pair{"dueDate", "Enter a due date"},
		"error":   pair{"dueDate", "Enter a valid value for due date"},
	},
	"notes": {
		"stringLengthTooLong": pair{"notes", "The note must be 1000 characters or fewer"},
		"isEmpty":             pair{"notes", "Enter a note explaining the issue"},
	},
	"taskCompletedNotes": {
		"stringLengthTooLong": pair{"notes", "The note must be 1000 characters or fewer"},
	},
	// other
	"otherImportantInformation": {
		"stringLengthTooLong": pair{"otherImportantInformation", "The other important information must be 1000 characters or fewer"},
	},
	"reportDueDate": {
		"error": pair{"reportDueDate", "Visit report due date - This must be a future date"},
	},
	"reportReceivedDate": {
		"invalid-lte": pair{"reportReceivedDate", "Report received date - This must be on or before " + currentDate},
	},
	"reportReviewDate": {
		"invalid-lte": pair{"reportReviewDate", "Report review date - This must be on or before " + currentDate},
	},
	"documentType": {
		"isEmpty": pair{"documentType", "Select the type of document"},
	},
	"documentDirection": {
		"isEmpty": pair{"documentDirection", "Select the document direction"},
	},
	"documentDate": {
		"isEmpty":     pair{"documentDate", "Enter a date"},
		"invalid-lte": pair{"documentDate", "The date must be today or in the past"},
	},
	"gcmIssueType": {
		"isEmpty": pair{"gcmIssueType", "Select an issue type"},
	},
	"caseRecNumber": {
		"checksumFailed": pair{"caseRecNumber", "Case number not recognised"},
	},
}

func RenameErrors(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	mappedErrors := sirius.ValidationErrors{}
	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if mapping, ok := validationMappings[fieldName][errorType]; ok {
				err[errorType] = mapping.v
				mappedErrors[mapping.k] = err
			} else {
				err[errorType] = errorMessage
				mappedErrors[fieldName] = err
			}
		}
	}
	return mappedErrors
}
