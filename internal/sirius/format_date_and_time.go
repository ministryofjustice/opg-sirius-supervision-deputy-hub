package sirius

import (
	"time"
)

const DateTimeFormat string = "2006-01-02T15:04:05+07:00"
const DateTimeDisplayFormat string = "2006-01-02"

func GenerateTimeForTest(year int, month time.Month, day, hour, min, seconds int) time.Time {
	newTime := time.Date(year, month, day, hour, min, seconds, 0, time.UTC)
	return newTime
}

func GetNullDate() time.Time {
	nullDate, _ := time.Parse("2006-01-02T15:04:05+00:00", "0001-01-01 00:00:00 +0000 UTC")
	return nullDate
}
