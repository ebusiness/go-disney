package utils

import (
	"time"
)

var (
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
	// jst, _ = time.LoadLocation("Asia/Tokyo")
)

// Now - time.Now() of Asia/Tokyo
func Now() time.Time {
	return time.Now().In(jst)
}

// DatetimeOfDate -
func DatetimeOfDate(datetime time.Time) time.Time {
	tokyoTime := datetime.In(jst)
	return time.Date(tokyoTime.Year(), tokyoTime.Month(), tokyoTime.Day(), 0, 0, 0, 0, jst)
}

// TokyoTime -
func TokyoTime(datetime time.Time) time.Time {
	return datetime.In(jst)
}
