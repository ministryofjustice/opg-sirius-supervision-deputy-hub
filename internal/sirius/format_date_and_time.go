package sirius

import "time"

func FormatDateAndTime(formatForDateTime string, dateString string, displayLayoutDateTime string) string {
	stringToDateTime, _ := time.Parse(formatForDateTime, dateString)
	dateTime := stringToDateTime.Format(displayLayoutDateTime)
	return dateTime
}
