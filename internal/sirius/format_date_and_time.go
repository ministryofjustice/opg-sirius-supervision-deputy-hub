package sirius

import (
	"time"
)

const DateTimeFormat string = "2006-01-02T15:04:05+00:00"
const DateTimeDisplayFormat string = "2006-01-02"

func FormatDateAndTime(formatForDateTime string, dateString string, displayLayoutDateTime string) string {
	if dateString == "" {
		return dateString
	}
	stringToDateTime, _ := time.Parse(formatForDateTime, dateString)
	dateTime := stringToDateTime.Local().Format(displayLayoutDateTime)
	return dateTime
}

func isDST() bool {
	return time.Now().Local() != time.Now()
}
