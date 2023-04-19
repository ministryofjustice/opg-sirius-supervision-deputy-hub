package sirius

import "time"

const DateTimeFormat string = "2006-01-02T15:04:05+00:00"
const DateTimeDisplayFormat string = "2006-01-02"

func FormatDateAndTime(formatForDateTime string, dateString string, displayLayoutDateTime string) string {
	loc, _ := time.LoadLocation("Europe/Dublin")

	if dateString == "" {
		return dateString
	}

	parsedTime, _ := time.Parse(formatForDateTime, dateString)
	dateTime := parsedTime.In(loc).Format(displayLayoutDateTime)

	return dateTime
}
