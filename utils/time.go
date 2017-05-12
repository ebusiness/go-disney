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
