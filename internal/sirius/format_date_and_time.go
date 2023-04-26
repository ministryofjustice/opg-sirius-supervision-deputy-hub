package sirius

import (
	"time"
)

const DateTimeFormat string = "2006-01-02T15:04:05+07:00"
const DateTimeDisplayFormat string = "2006-01-02"

func GenerateTimeForTest(year int, month time.Month, day, hour, min, seconds int) time.Time {
	//location := time.FixedZone("", 0)
	newTime := time.Date(year, month, day, hour, min, seconds, 0, time.UTC)
	return newTime
}
