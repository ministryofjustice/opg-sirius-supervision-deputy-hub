package sirius

import (
	"time"
)

const DateTimeFormat string = "2006-01-02T15:04:05+07:00"
const DateTimeDisplayFormat string = "2006-01-02"

func FormatDateTimeStringIntoDateTime(formatForDateTime string, dateString string) time.Time {
	location, _ := time.LoadLocation("UTC")

	date, _ := time.Parse(formatForDateTime, dateString)
	return date.In(location)
}
