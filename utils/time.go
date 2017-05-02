package utils

import (
	"time"
)

var (
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
)

// Now - time.Now() of Asia/Tokyo
func Now() time.Time {
	return time.Now().In(jst)
}