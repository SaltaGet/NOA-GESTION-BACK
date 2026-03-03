package utils

import "time"

func GetFirstAndLastDayTwoMonthsAgo(month int) (time.Time, time.Time) {
	target := time.Now().AddDate(0, -month, 0)

	first := time.Date(target.Year(), target.Month(), 1, 0, 0, 0, 0, target.Location())

	last := first.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return first, last
}