package sirius

import (
	"regexp"
	"time"
)

const IsoDateTimeZone string = "2006-01-02T15:04:05+00:00"
const IsoDateTime string = "2006-01-02 15:04:05"
const IsoDate string = "2006-01-02"
const SiriusDate string = "02/01/2006"
const SiriusDateTime string = "02/01/2006 15:04:05"

func FormatDateTime(currentFormat string, dateString string, displayFormat string) string {
	if dateString == "" {
		return dateString
	}
	stringToDateTime, _ := time.Parse(currentFormat, dateString)
	dateTime := stringToDateTime.Local().Format(displayFormat)
	return dateTime
}

// PermissiveFormatIsoDateTime allows transition from `IsoDateTime` to `IsoDateTimeZone` without display errors for users.
// Once the transition is completed, this should be deleted and `FormatDateTime(IsoDateTimeZone, ...` used.
func PermissiveFormatIsoDateTime(dateString string, displayFormat string) string {
	if dateString == "" {
		return dateString
	}

	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}$`)
	if re.MatchString(dateString) {
		stringToDateTime, _ := time.Parse(IsoDateTimeZone, dateString)
		return stringToDateTime.Local().Format(displayFormat)
	}

	stringToDateTime, _ := time.Parse(IsoDateTime, dateString)
	return stringToDateTime.Local().Format(displayFormat)
}
