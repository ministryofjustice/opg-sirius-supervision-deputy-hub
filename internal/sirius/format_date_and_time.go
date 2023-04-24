package sirius

import (
	"time"
)

const DateTimeFormat string = "2006-01-02T15:04:05+07:00"
const DateTimeDisplayFormat string = "2006-01-02"

func FormatDateTimeStringIntoDateTime(formatForDateTime string, dateString string) time.Time {
	stringToDateTime, _ := time.Parse(formatForDateTime, dateString)
	return stringToDateTime
}
