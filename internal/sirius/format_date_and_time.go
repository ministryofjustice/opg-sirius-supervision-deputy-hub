package sirius

import "time"

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
